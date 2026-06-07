"""REST API endpoints for the www.fmind.dev portfolio application."""

from fastapi import APIRouter

from app.data.certifications import BADGES, SPECIALIZATIONS
from app.data.experience import EXPERIENCES
from app.data.projects import OPEN_SOURCE, YOUTUBE_SERIES

router = APIRouter()


@router.get("/api/certifications")
async def get_certifications():
    """Returns lists of certification badges and specializations."""
    return {"badges": BADGES, "specializations": SPECIALIZATIONS}


@router.get("/api/projects")
async def get_projects():
    """Returns lists of open-source projects and YouTube series."""
    return {"open_source": OPEN_SOURCE, "youtube_series": YOUTUBE_SERIES}


@router.get("/api/experience")
async def get_experience():
    """Returns list of job experience items."""
    return {"experiences": EXPERIENCES}
