# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.26 AS builder
ARG TARGETARCH
ARG TARGETOS
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go tool templ generate && \
  CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
  go build -trimpath -ldflags="-s -w" -o www-fmind-dev ./cmd/www-fmind-dev

FROM gcr.io/distroless/static-debian13:nonroot AS runner
ENV ENVIRONMENT=production
USER nonroot:nonroot
WORKDIR /app
EXPOSE 8080
COPY --from=builder /app/www-fmind-dev ./www-fmind-dev
ENTRYPOINT ["./www-fmind-dev"]
