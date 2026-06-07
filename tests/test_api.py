from fastapi.testclient import TestClient

from app.main import app

client = TestClient(app)


def test_get_certifications():
    """Tests GET /api/certifications schema and content."""
    response = client.get("/api/certifications")
    assert response.status_code == 200
    data = response.json()
    assert "badges" in data
    assert "specializations" in data
    assert len(data["badges"]) > 0
    assert len(data["specializations"]) > 0


def test_get_projects():
    """Tests GET /api/projects schema and content."""
    response = client.get("/api/projects")
    assert response.status_code == 200
    data = response.json()
    assert "open_source" in data
    assert "youtube_series" in data
    assert len(data["open_source"]) > 0
    assert len(data["youtube_series"]) > 0


def test_get_experience():
    """Tests GET /api/experience schema and content."""
    response = client.get("/api/experience")
    assert response.status_code == 200
    data = response.json()
    assert "experiences" in data
    assert len(data["experiences"]) > 0


def test_healthz():
    """Tests GET /healthz response status and details."""
    response = client.get("/healthz")
    assert response.status_code == 200
    data = response.json()
    assert "status" in data
    assert data["status"] in ("ok", "warning", "error")
