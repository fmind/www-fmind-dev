"""Middleware components for the www.fmind.dev FastAPI application.

This module defines standard HTTP middleware functions for structured logging,
GCP trace context injection, secure browser headers (CSP, HSTS), and naked domain
redirection.
"""

from __future__ import annotations

import os
from typing import TYPE_CHECKING

import structlog

from app.config import settings

if TYPE_CHECKING:
    from fastapi import Request

logger = structlog.get_logger()


async def trace_and_log_middleware(request: Request, call_next):
    """Middleware to inject GCP trace IDs and log incoming requests."""
    structlog.contextvars.clear_contextvars()

    # Extract GCP trace context if present
    trace_header = request.headers.get("x-cloud-trace-context")
    if trace_header:
        trace_id = trace_header.split("/")[0]
        project_id = os.environ.get("GOOGLE_CLOUD_PROJECT") or "www-fmind-dev"
        structlog.contextvars.bind_contextvars(
            **{"logging.googleapis.com/trace": f"projects/{project_id}/traces/{trace_id}"}
        )

    logger.info("request_started", method=request.method, path=request.url.path)
    try:
        response = await call_next(request)
        logger.info("request_finished", status_code=response.status_code)
        return response
    except Exception as e:
        logger.exception("request_failed", error=str(e))
        raise e


async def security_headers_middleware(request: Request, call_next):
    """Enforces secure browser headers (CSP, HSTS, X-Frame-Options)."""
    response = await call_next(request)
    response.headers["X-Frame-Options"] = "DENY"
    response.headers["X-Content-Type-Options"] = "nosniff"
    response.headers["X-XSS-Protection"] = "1; mode=block"

    # CSP supporting local assets and Alpine.js inline directives
    response.headers["Content-Security-Policy"] = (
        "default-src 'self'; "
        "script-src 'self' 'unsafe-inline' 'unsafe-eval'; "
        "style-src 'self' 'unsafe-inline'; "
        "img-src 'self' data:; "
        "font-src 'self' data:; "
        "frame-src 'self';"
    )

    # Enforce HTTPS HSTS in production env
    if settings.is_production:
        response.headers["Strict-Transport-Security"] = "max-age=63072000; includeSubDomains; preload"

    return response


async def redirect_naked_domain_middleware(request: Request, call_next):
    """Redirects fmind.dev to www.fmind.dev."""
    host = request.headers.get("host", "")
    if host.split(":")[0] == "fmind.dev":
        from fastapi.responses import RedirectResponse

        redirect_url = str(request.url.replace(scheme="https", netloc="www.fmind.dev"))
        return RedirectResponse(url=redirect_url, status_code=301)
    return await call_next(request)
