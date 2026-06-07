# https://docs.docker.com/reference/dockerfile/

# Stage 1: Build python virtual environment
FROM python:3.14-slim AS builder

# Install uv
COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

# Set working directory
WORKDIR /app

# Copy python dependencies definitions
COPY pyproject.toml uv.lock .python-version ./

# Create virtualenv and sync dependencies
RUN --mount=type=cache,target=/root/.cache/uv \
  uv sync --frozen --no-install-project

# Stage 2: Final runner image
FROM python:3.14-slim AS runner

# Create a non-root system user and group for security
RUN groupadd --system -g 10001 app && \
  useradd --system -u 10001 -g app --create-home --shell /sbin/nologin app

WORKDIR /app

# Copy virtualenv from builder stage
COPY --chown=app:app --from=builder /app/.venv /app/.venv

# Copy static assets and application code directly
COPY --chown=app:app static ./static
COPY --chown=app:app app ./app

# Set environment variables
ENV PATH="/app/.venv/bin:$PATH" \
  ENV=prod \
  PYTHONUNBUFFERED=1 \
  PYTHONDONTWRITEBYTECODE=1

# Run as non-root user
USER app

# Expose default port
EXPOSE 8080

# Execute command using granian
CMD ["/bin/sh", "-c", "exec granian --interface asgi --host 0.0.0.0 --port ${PORT:-8080} --workers 1 --runtime-threads 4 app.main:app"]
