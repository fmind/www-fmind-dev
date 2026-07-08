# syntax=docker/dockerfile:1
# Stage 1: Build the Go binary and compile Tailwind CSS.
# Debian 13 "trixie" (glibc) builder: the tailwindcss-extra CLI is glibc-linked, and trixie
# matches the distroless-debian13 runtime below (same glibc line). Runs natively on
# $BUILDPLATFORM and cross-compiles Go to the target — fast (no QEMU) and multi-arch.
FROM --platform=$BUILDPLATFORM golang:1.26-trixie AS builder

# Provided automatically by BuildKit.
ARG BUILDARCH
ARG TARGETOS
ARG TARGETARCH

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# hadolint ignore=DL3008
RUN apt-get update \
    && apt-get install -y --no-install-recommends curl ca-certificates git \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Download + verify the standalone tailwindcss-extra binary for the BUILD arch
# (Tailwind runs during the build; the CSS it emits is architecture-independent).
# v2.9.1 bundles DaisyUI 5.6.13, matching the mise-managed CLI used in local dev
# (its `include:` config is required to emit the slimmed stylesheet).
RUN mkdir -p bin && \
    case "${BUILDARCH}" in \
      amd64) TW_ARCH=x64;   TW_SHA=c3486e01633d39225aa1084f4a0825189145459ec628be5df19b0d51012665d4 ;; \
      arm64) TW_ARCH=arm64; TW_SHA=49e87dde2d2c82839b6c4f20c2fcab7e7a5c747243a4fa27f854caddb0a03278 ;; \
      *) echo "unsupported build arch: ${BUILDARCH}" && exit 1 ;; \
    esac && \
    curl -sSL "https://github.com/dobicinaitis/tailwind-cli-extra/releases/download/v2.9.1/tailwindcss-extra-linux-${TW_ARCH}" -o bin/tailwindcss-extra && \
    echo "${TW_SHA}  bin/tailwindcss-extra" | sha256sum -c - && \
    chmod +x bin/tailwindcss-extra

# Dependency layer (cached across source changes).
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

# Compile Templ templates.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go tool templ generate

# Compile + minify Tailwind CSS.
RUN ./bin/tailwindcss-extra -i static/css/input.css -o static/dist/styles.css --minify

# Cross-compile the static Go binary for the target platform.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o www-fmind-dev ./cmd/www-fmind-dev

# Stage 2: Development runner (Go toolchain + air + tailwind watch). Used by `skaffold dev`.
FROM builder AS dev-runner
ENV ENVIRONMENT=development
EXPOSE 8080
CMD ["bash", "-c", "./bin/tailwindcss-extra -i static/css/input.css -o static/dist/styles.css --watch & go tool air"]

# Stage 3: Production runner — distroless static, nonroot (uid 65532), shell-less.
# Kept LAST so a plain `docker build` (no --target) defaults to the hardened image.
# debian13 is the current distroless line (debian12 is superseded). Renovate pins the tag;
# for full supply-chain reproducibility consider pinning by @sha256 digest.
FROM gcr.io/distroless/static-debian13:nonroot AS runner
WORKDIR /app
COPY --from=builder /app/www-fmind-dev ./www-fmind-dev
ENV ENVIRONMENT=production
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["./www-fmind-dev"]
