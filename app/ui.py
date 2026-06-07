from fasthtml.common import *

from app.config import settings
from app.layout import get_headers, page_layout
from app.sections.about import about_section
from app.sections.certifications import certifications_section
from app.sections.hero import hero_section
from app.sections.projects import projects_section
from app.sections.publications import publications_section
from app.sections.services import services_section
from app.sections.work_experience import work_experience_section


def _get_session_secret() -> str:
    """Gets the session signing secret from the environment.

    In production, a SESSION_SECRET environment variable must be set.
    In local development, it falls back to a development-only default.
    """
    if settings.is_production and (
        settings.session_secret == "dev-secret-key-change-me-in-production"  # noqa: S105
        or not settings.session_secret.strip()
    ):
        raise RuntimeError("SESSION_SECRET environment variable is required in production.")

    return settings.session_secret


rt_app = FastHTML(
    title="www.fmind.dev — Médéric Hurier | AI/ML Architect & Engineer",
    hdrs=get_headers(),
    default_hdrs=False,
    htmlkw={"lang": "en", "class": "scroll-smooth"},
    secret_key=_get_session_secret(),
)
rt = rt_app.route


@rt("/")
def home():
    """Main route handler for www.fmind.dev portfolio home page."""
    content = Div(
        hero_section(),
        about_section(),
        work_experience_section(),
        certifications_section(),
        publications_section(),
        projects_section(),
        services_section(),
    )
    return page_layout(content)
