# AI Agent Guidance & conventions

This document outlines the technical conventions, tooling, commands, and layout of this repository for AI coding assistants.

## Tooling & Commands

*   **Virtual Environment**: Managed exclusively by `uv`. All python commands must run within `.venv` (e.g. `uv run python`).
*   **Task Automation**: Managed via `just`. Run `just` to see the list of available recipes.
*   **Formatting & Linting**: Managed by `ruff`. Run `just format` to format code, or `just check` to run code verification checks.
*   **Static Type Checking**: Managed via `ty` (wrapping pyright/mypy). Run `just check` to verify types.
*   **Git Hooks**: Managed via `lefthook` for pre-commit lint/format/type-checks and pre-push test validation.
*   **Interactive Styles**: Compiled locally via Tailwind CSS CLI. Run `just css` to compile `static/dist/styles.css`.

## Conventions

*   **100% Self-Hosted**: No remote script or style sheets are allowed in page layouts. All assets (fonts, images, js libraries like HTMX and Alpine.js) must be hosted under `/static/` and served locally.
*   **Python Code Style**: Enforce strict snake_case naming for functions, methods, and variables. Write Google-style docstrings for all public modules, functions, and classes.
*   **Dynamic Attributes**: Unpacking dictionaries for Alpine directives (e.g. `**{"@click": "..."}`) is standard for Starlette/FastHTML components to avoid type-check errors.

## Layout

*   `app/` — FastAPI REST endpoints, FastHTML UI routes, models, and local static data files.
*   `infra/` — Terraform configurations for Cloud Run deployment, APIs, IAM, and Cloud Monitoring.
*   `static/` — Compiled styles, self-hosted JS libraries (HTMX/Alpine), fonts, and image assets.
*   `tests/` — Pytest unit and integration test suite for REST APIs and HTML page delivery.
*   `.env` — Local development environment variables setup file.
*   `.python-version` — Pinned python runtime version (3.14).
*   `Dockerfile` — Production-ready multi-stage build starting with python:3.14-slim.
*   `.github/workflows/deploy.yml` — GitHub Actions CI/CD deployment workflow.
*   `justfile` — Developer task recipes automation file.
*   `lefthook.yml` — Pre-commit and pre-push Git hooks definition.
*   `pyproject.toml` — uv packages, project configuration, and QA checks settings.
*   `AGENTS.md` (this file) — AI agent conventions and repository technical layout.
*   `README.md` — Project overview and developer setup guide.
