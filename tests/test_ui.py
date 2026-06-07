from fastapi.testclient import TestClient

from app.main import app

client = TestClient(app)


def test_home_page_delivery():
    """Tests that the root route renders the portfolio layout, sections, SEO tags, and local assets correctly."""
    response = client.get("/")
    assert response.status_code == 200
    html = response.text

    # 1. Critical SEO and Metadata Tags
    assert "Médéric Hurier | AI/ML Architect &amp; Engineer" in html
    assert 'name="description"' in html
    assert 'content="Médéric Hurier (Fmind)"' in html
    assert 'name="viewport"' in html

    # 2. 100% Self-Hosted Local Asset References and Linked Styling
    assert 'href="/static/dist/styles.css?v=' in html
    assert 'src="/static/js/htmx.min.js"' in html
    assert 'src="/static/js/alpine.min.js"' in html

    # 3. Presence of all main portfolio sections
    assert 'id="about"' in html
    assert 'id="work-experience"' in html
    assert 'id="certifications"' in html
    assert 'id="publications"' in html
    assert 'id="projects"' in html
    assert 'id="services"' in html

    # 4. JSON-LD Structured Data Schema
    assert 'type="application/ld+json"' in html


def test_root_static_files():
    """Tests that root-level static metadata files are successfully served."""
    routes = [
        "/robots.txt",
        "/sitemap.xml",
        "/favicon.ico",
        "/llms.txt",
        "/humans.txt",
        "/site.webmanifest",
        "/security.txt",
        "/.well-known/security.txt",
    ]
    for route in routes:
        response = client.get(route)
        assert response.status_code == 200, f"Failed to serve {route}"


def test_naked_domain_redirect():
    """Tests that requests to the naked domain fmind.dev redirect to www.fmind.dev."""
    response = client.get("/", headers={"Host": "fmind.dev"}, follow_redirects=False)
    assert response.status_code == 301
    assert response.headers["location"] == "https://www.fmind.dev/"

    # Test that path and query params are preserved
    response = client.get("/healthz?test=1", headers={"Host": "fmind.dev"}, follow_redirects=False)
    assert response.status_code == 301
    assert response.headers["location"] == "https://www.fmind.dev/healthz?test=1"

    # Test that www.fmind.dev is not redirected
    response = client.get("/", headers={"Host": "www.fmind.dev"})
    assert response.status_code == 200

    # Test that localhost is not redirected
    response = client.get("/", headers={"Host": "localhost:8080"})
    assert response.status_code == 200
