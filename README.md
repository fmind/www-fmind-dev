# www.fmind.dev — Médéric Hurier Portfolio

The serverless portfolio website of Médéric Hurier (Fmind), freelance AI/ML Architect & Engineer, built using FastAPI + FastHTML + Tailwind CSS/DaisyUI + HTMX + Alpine.js and compiled locally with Tailwind CSS v4.

## Prerequisites

*   **Python**: Version 3.14+ (CPython)
*   **uv**: Dependency and workspace manager
*   **just**: Task runner

## Local Development Setup

1.  **Clone the Repository**:
    ```bash
    git clone <repository_url> www-fmind-dev
    cd www-fmind-dev
    ```

2.  **Configure Environment Variables**:
    Copy the sample environment template to `.env`:
    ```bash
    cp .env.sample .env
    ```
    Local development defaults are pre-configured, but in production, ensure `SESSION_SECRET` is set to a secure secret.

3.  **Install Dependencies**:
    ```bash
    just install
    ```
    This syncs the python virtualenv and registers the lefthook git hooks.

4.  **Run Development Server**:
    ```bash
    just dev
    ```
    This compiles the Tailwind CSS assets and starts the reload-enabled ASGI development server. Access the site at `http://localhost:8080`.

## Testing & Checks

*   **Run All Quality Checks**:
    ```bash
    just check
    ```
    This verifies formatting, lint errors, and runs the `ty` static type checker.

*   **Format Code**:
    ```bash
    just format
    ```
    This sorts imports and reformats all files using ruff.

*   **Run Tests**:
    ```bash
    just test
    ```
    This executes the pytest test suite verifying HTML delivery, REST API schemas, and health checks.

## Production & Deployment

*   **Build & Run Locally (Production Mode)**:
    ```bash
    just prod
    ```
    This builds the final stylesheet and launches the high-performance Granian ASGI server.

*   **Docker Build**:
    ```bash
    docker build -t www-fmind-dev .
    ```
    The multi-stage build uses `python:3.14-slim` and completely separates Node/Tailwind build stages to keep the runner image minimal and secure.
