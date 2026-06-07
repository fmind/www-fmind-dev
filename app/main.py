"""Main entry point for the www.fmind.dev FastAPI and FastHTML application.

This module composes the application by setting up logging, defining global
middlewares, mounting static assets, serving root-level metadata files, and
mounting the REST API router and FastHTML UI sub-app.
"""

import httpx
import structlog
from fastapi import FastAPI, Response
from fastapi.middleware.gzip import GZipMiddleware
from fastapi.responses import FileResponse
from fastapi.staticfiles import StaticFiles

from app.api import router as api_router
from app.logging_config import setup_logging
from app.middleware import (
    redirect_naked_domain_middleware,
    security_headers_middleware,
    trace_and_log_middleware,
)
from app.ui import rt_app

# Initialize structured logging
setup_logging()
logger = structlog.get_logger()


class CachedStaticFiles(StaticFiles):
    """Static files hosting with long-term caching headers."""

    def file_response(self, full_path, stat_result, scope, status_code=200):
        response = super().file_response(full_path, stat_result, scope, status_code)
        response.headers["Cache-Control"] = "public, max-age=31536000, immutable"
        return response


# Parent FastAPI app
app = FastAPI(title="www.fmind.dev API", version="1.0.0")

# Register middlewares (last registered executes first on incoming requests)
app.add_middleware(GZipMiddleware, minimum_size=1000)
app.middleware("http")(trace_and_log_middleware)
app.middleware("http")(security_headers_middleware)
app.middleware("http")(redirect_naked_domain_middleware)


# --- Root Static Files & Health Routes ---


@app.get("/health")
@app.get("/healthz")
async def health_check(response: Response):
    """Lightweight health check verifying static files and network connectivity."""
    status = "ok"
    details = {}

    # 1. Static files check
    from pathlib import Path

    css_path = Path("static/dist/styles.css")
    details["static_files_ok"] = css_path.exists()
    if not details["static_files_ok"]:
        status = "error"

    # 2. Downstream internet connectivity check
    try:
        async with httpx.AsyncClient(timeout=2.0) as client:
            res = await client.get("https://www.google.com")
            details["network_ok"] = res.status_code == 200
            if not details["network_ok"]:
                status = "warning"
    except Exception as e:
        details["network_ok"] = False
        details["network_error"] = str(e)
        status = "warning"

    if status == "error":
        response.status_code = 503

    return {"status": status, "details": details}


@app.get("/robots.txt", include_in_schema=False)
async def serve_robots():
    return FileResponse("static/robots.txt")


@app.get("/sitemap.xml", include_in_schema=False)
async def serve_sitemap():
    return FileResponse("static/sitemap.xml")


@app.get("/favicon.ico", include_in_schema=False)
async def serve_favicon():
    return FileResponse("static/favicon.ico")


@app.get("/llms.txt", include_in_schema=False)
async def serve_llms():
    return FileResponse("static/llms.txt")


@app.get("/humans.txt", include_in_schema=False)
async def serve_humans():
    return FileResponse("static/humans.txt")


@app.get("/site.webmanifest", include_in_schema=False)
async def serve_manifest():
    return FileResponse("static/site.webmanifest")


@app.get("/security.txt", include_in_schema=False)
@app.get("/.well-known/security.txt", include_in_schema=False)
async def serve_security():
    return FileResponse("static/security.txt")


# Mount Static Files first
app.mount("/static", CachedStaticFiles(directory="static"), name="static")

# Include REST API endpoints
app.include_router(api_router)

# Mount FastHTML at the root (catches all HTML requests)
app.mount("/", rt_app)

if __name__ == "__main__":
    import os

    import uvicorn

    port = int(os.environ.get("PORT", 8080))
    uvicorn.run("app.main:app", host="0.0.0.0", port=port, reload=True)  # noqa: S104
