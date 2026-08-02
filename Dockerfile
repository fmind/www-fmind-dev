# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.26.5 AS css
ARG BUILDARCH
ADD --chmod=0755 --checksum=sha256:a6f21b9ae0970e927dd374e6a25a17936a11f5d4133cd825b0ca6492cc9b1b31 \
  https://github.com/dobicinaitis/tailwind-cli-extra/releases/download/v2.10.8/tailwindcss-extra-linux-x64 \
  /usr/local/lib/tailwindcss-extra-amd64
ADD --chmod=0755 --checksum=sha256:1b60d91961df32b43590fbdfe63fd26f4e2db585d7a2beb56256ca27b6b18cdd \
  https://github.com/dobicinaitis/tailwind-cli-extra/releases/download/v2.10.8/tailwindcss-extra-linux-arm64 \
  /usr/local/lib/tailwindcss-extra-arm64
WORKDIR /app
COPY static/css/input.css static/css/input.css
COPY templates templates
RUN /usr/local/lib/tailwindcss-extra-${BUILDARCH} -i static/css/input.css -o static/dist/styles.css --minify

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
