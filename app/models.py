"""Data models for the portfolio website."""

from typing import Literal

from pydantic import BaseModel


class ExpertiseCard(BaseModel):
    """Core expertise & skills card model."""

    title: str
    emoji: str
    gradient: str
    description: str


class CertificationBadge(BaseModel):
    """Credential badge model."""

    url: str
    logo: str
    title: str
    issuer: str
    cert_id: str
    active: bool


class CertificationEntry(BaseModel):
    """Specialization and foundation certificate model."""

    url: str
    logo: str
    title: str
    issuer_details: str


class Project(BaseModel):
    """Open-source project model."""

    title: str
    href: str
    repo: str | None = None
    description: str
    type: Literal["github", "youtube"]


class Playlist(BaseModel):
    """YouTube series model."""

    title: str
    url: str
    description: str
    cta: str


class ResearchPaper(BaseModel):
    """Academic research paper model."""

    title: str
    url: str
    venue: str
    code: str
    code_label: str


class CuratedPost(BaseModel):
    """Featured post/publication model."""

    title: str
    url: str


class Service(BaseModel):
    """Service offering model."""

    icon: str
    title: str
    description: str
    badge: str
    badge_class: str
    cta_text: str
    cta_url: str


class WorkExperience(BaseModel):
    """Job history timeline entry model."""

    company: str
    logo: str
    title: str
    color: str
    hex: str
    description: str
    tags: list[str]
