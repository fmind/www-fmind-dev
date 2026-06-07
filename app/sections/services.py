from fasthtml.common import *

from app.data.services import SERVICES


def services_section():
    """Renders the portfolio Services section."""
    # Build Services Cards Grid
    service_cards = [
        Div(
            Div(
                # Large Icon
                Span(
                    service.icon,
                    cls="text-6xl mb-4 inline-block group-hover:scale-110 transition-transform duration-300 select-none",
                ),
                H3(service.title, cls="card-title text-2xl font-heading font-bold text-base-content select-text"),
                P(
                    service.description,
                    cls="text-base-content/70 mt-3 text-sm md:text-base leading-relaxed select-text",
                ),
                # Status Badge
                Div(service.badge, cls=f"badge gap-1 mt-6 font-semibold py-3 px-4 select-none {service.badge_class}"),
                # CTA Button
                Div(
                    A(
                        service.cta_text,
                        href=service.cta_url,
                        target="_blank",
                        rel="noopener noreferrer",
                        cls="btn btn-primary bg-primary hover:bg-primary/90 text-white border-none rounded-xl shadow-md hover:scale-105 transition-transform hover:shadow-lg w-full max-w-[220px] py-3 text-center font-bold",
                    ),
                    cls="card-actions mt-8 w-full justify-center",
                ),
                cls="card-body p-8 items-center text-center",
            ),
            cls="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 border-2 border-primary group rounded-2xl overflow-hidden",
        )
        for service in SERVICES
    ]

    return Section(
        Div(
            # Section Header
            Div(
                H2(
                    "My ",
                    Span("Services", cls="text-primary"),
                    cls="text-4xl md:text-5xl font-heading font-bold mb-4 text-primary",
                ),
                P(
                    "Leverage deep AI/ML expertise for your next enterprise-level initiative.",
                    cls="text-xl text-base-content/70 max-w-3xl mx-auto leading-relaxed text-balance",
                ),
                cls="text-center mb-16",
            ),
            # Services Cards Grid
            Div(
                *service_cards,
                cls="grid grid-cols-1 md:grid-cols-2 gap-8 max-w-4xl mx-auto",
            ),
            cls="container mx-auto",
        ),
        id="services",
        cls="py-24 bg-base-200 overflow-hidden relative px-4",
    )
