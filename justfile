# https://just.systems/man/en/

set dotenv-load := true
export UV_ENV_FILE := ".env"

image := "www-fmind-dev"
package := "app"
port := env_var_or_default("PORT", "8080")
project_id := "www-fmind-dev"


# show all available recipes
default:
    @just --list

# run all static checkers
check:
    terraform fmt -check -diff infra/
    uv run ruff format --check .
    uv run ruff check .
    uv run ty check
    uv lock --check

# clean caches and temps
clean:
    find . -type f -name '*.py[co]' -delete
    find . -type d -name __pycache__ -exec rm -rf {} +
    rm -rf .pytest_cache .ruff_cache .coverage htmlcov docs

# run and open coverage
coverage:
    uv run pytest --cov --cov-report=html
    xdg-open htmlcov/index.html || open htmlcov/index.html

# compile Tailwind CSS
css:
    uv run tailwindcss-extra -i static/css/input.css -o static/dist/styles.css --minify

# run local dev server
dev: css
    uv run python -m app.main

# build and run docker container locally
docker *args="": css
    docker build -t {{ image }} .
    docker run -it --rm -e ENV -e SESSION_SECRET -e PORT={{ port }} -p {{ port }}:{{ port }} {{ args }} {{ image }}

# sort imports and auto-format code
format:
    terraform fmt infra/
    uv run ruff check --select=I --fix .
    uv run ruff format .

# install dependencies and hooks
install:
    uv sync
    uv run lefthook install

# run test suite with pytest
test *args="":
    uv run pytest {{ args }}

# run terraform commands
tf *args="":
    terraform -chdir=infra {{ args }}

# upgrade all packages
upgrade:
    uv sync --upgrade

# vendor static assets
vendor:
    uv run python scripts/vendor.py

# generate and push a new production SESSION_SECRET version to GCP Secret Manager
[confirm("Are you sure you want to generate and push a new production SESSION_SECRET version to GCP?")]
secrets:
    openssl rand -hex 32 | gcloud secrets versions add session-secret --project={{ project_id }} --data-file=-

