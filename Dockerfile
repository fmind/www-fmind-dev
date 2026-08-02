# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.26.5 AS css
ARG BUILDARCH
ADD --chmod=0755 --checksum=sha256:6990c382dac92f392231e1e75863dd1cd215323038348fdd69e61d3fcdc8aa4c \
  https://github.com/dobicinaitis/tailwind-cli-extra/releases/download/v2.10.5/tailwindcss-extra-linux-x64 \
  /usr/local/lib/tailwindcss-extra-amd64
ADD --chmod=0755 --checksum=sha256:13cba246aec769d4e4dffacc252888590f869bea0089c791e839586ae6577372 \
  https://github.com/dobicinaitis/tailwind-cli-extra/releases/download/v2.10.5/tailwindcss-extra-linux-arm64 \
  /usr/local/lib/tailwindcss-extra-arm64
WORKDIR /app
COPY static/css/input.css static/css/input.css
COPY templates templates
RUN case "${BUILDARCH}" in \
      amd64|arm64) /usr/local/lib/tailwindcss-extra-${BUILDARCH} -i static/css/input.css -o static/dist/styles.css --minify ;; \
      *) echo "unsupported build architecture: ${BUILDARCH}" >&2; exit 1 ;; \
    esac

FROM --platform=$BUILDPLATFORM golang:1.26.5 AS builder
ARG TARGETARCH
ARG TARGETOS
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
COPY --from=css /app/static/dist/styles.css static/dist/styles.css
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go tool templ generate && \
  CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
  go build -trimpath -ldflags="-s -w" -o www-fmind-dev ./cmd/www-fmind-dev

FROM gcr.io/distroless/static-debian13:nonroot AS runner
ENV ENVIRONMENT=production
# Distroless defines nonroot as numeric 65532; numeric ownership remains
# resolvable even when a container runtime does not load /etc/passwd.
USER 65532:65532
WORKDIR /app
EXPOSE 8080
COPY --from=builder /app/www-fmind-dev ./www-fmind-dev
ENTRYPOINT ["./www-fmind-dev"]
