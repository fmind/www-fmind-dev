from datetime import datetime

from app.models import Service

SERVICES = [
    Service(
        icon="🏢",
        title="Freelancing",
        description=(
            "Hire me as a freelance AI/ML Architect & Engineer to design, build, and scale your AI initiatives — "
            "from AI Agents to MLOps."
        ),
        badge=f"⏳ Available from {datetime.now().year}",
        badge_class="badge-warning text-warning-content",
        cta_text="✉️ Contact Me",
        cta_url="mailto:contact@fmind.dev",
    ),
    Service(
        icon="🎓",
        title="Mentoring",
        description=(
            "Book a paid 1-hour session to discuss your projects: upskilling, career mentoring, "
            "architecture review, brainstorming, and more."
        ),
        badge="💰 Paid session — 1 hour",
        badge_class="badge-info text-info-content",
        cta_text="📅 Book a Session",
        cta_url=(
            "https://calendar.google.com/calendar/u/0/appointments/schedules/"
            "AcZssZ2ye3X9589PA2xmbV73Iz5J_NbFig6nN651vn6UuYAC-Cs5vBxnQ2L5db9UnAeXmUBQSW1MOobd"
        ),
    ),
]
