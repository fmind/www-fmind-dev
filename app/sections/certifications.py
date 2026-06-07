from fasthtml.common import *

from app.data.certifications import BADGES, SPECIALIZATIONS
from app.icons import icon


def certifications_section():
    """Renders the portfolio Certifications section."""
    # Build Badge Links
    badge_links = [
        A(
            # Logo Container
            Div(
                Img(
                    src=f"/static/img/certifications/{badge.logo}",
                    alt=badge.issuer,
                    cls="w-11 h-11 object-contain",
                    width="44",
                    height="44",
                ),
                cls="bg-base-200 p-2.5 rounded-xl group-hover:bg-base-300 transition-colors shrink-0 flex items-center justify-center w-16 h-16 bg-slate-50 border border-base-200",
            ),
            # Content Info
            Div(
                H3(
                    badge.title,
                    cls=(
                        f"font-bold text-base md:text-md leading-tight line-clamp-2 transition-colors select-text "
                        f"{'text-primary' if badge.active else 'text-base-content group-hover:text-accent'}"
                    ),
                ),
                P(
                    f"{badge.issuer} • {badge.cert_id}",
                    cls="text-xs text-base-content/60 mt-1.5 truncate select-text",
                    title=f"{badge.issuer} • {badge.cert_id}",
                ),
                cls="flex-1 min-w-0",
            ),
            # Link Indicator Icon
            icon(
                "external_link",
                cls=(
                    f"w-4 h-4 transition-opacity shrink-0 "
                    f"{'text-primary opacity-50 group-hover:opacity-100' if badge.active else 'text-accent opacity-0 group-hover:opacity-100'}"
                ),
            ),
            href=badge.url,
            target="_blank",
            rel="noopener noreferrer",
            cls=(
                f"group relative flex items-center gap-3.5 p-4 bg-base-100 rounded-2xl border-2 transition-all "
                f"duration-300 w-full "
                f"{'border-primary shadow-lg hover:shadow-primary/20 hover:-translate-y-1' if badge.active else 'border-base-content/20 shadow-sm hover:border-base-content/40 hover:-translate-y-0.5'}"
            ),
        )
        for badge in BADGES
    ]

    # Build Specializations Entries
    spec_entries = [
        A(
            # Mini Logo
            Img(
                src=f"/static/img/certifications/{entry.logo}",
                alt=entry.title,
                cls="w-10 h-10 object-contain shrink-0",
                width="40",
                height="40",
            ),
            # Entry Info
            Div(
                H4(
                    entry.title,
                    cls="text-sm md:text-base font-semibold text-base-content group-hover:text-primary transition-colors line-clamp-2 select-text",
                ),
                P(entry.issuer_details, cls="text-xs md:text-sm text-base-content/70 mt-0.5 select-text"),
                cls="min-w-0 flex-1",
            ),
            # Indicator
            icon(
                "external_link",
                cls="w-4 h-4 ml-auto text-base-content/30 group-hover:text-primary transition-colors opacity-0 group-hover:opacity-100 shrink-0",
            ),
            href=entry.url,
            target="_blank",
            rel="noopener noreferrer",
            cls="group flex items-center gap-4 p-4 rounded-xl hover:bg-base-200 border border-transparent hover:border-base-300 transition-all cursor-pointer",
        )
        for entry in SPECIALIZATIONS
    ]

    return Section(
        Div(
            # Section Header
            Div(
                H2("Certifications", cls="text-4xl md:text-5xl font-heading font-bold mb-4 text-primary"),
                P(
                    "Validated capabilities from leading cloud providers and academic networks.",
                    cls="text-xl text-base-content/70 max-w-3xl mx-auto leading-relaxed text-balance",
                ),
                cls="text-center mb-16",
            ),
            # Credentials Badges Grid
            Div(
                *badge_links,
                cls="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 justify-items-center max-w-6xl mx-auto mb-16",
            ),
            # Specializations & Foundations Divider
            Div(
                Div(
                    Div(cls="w-full border-t border-base-300"),
                    cls="absolute inset-0 flex items-center",
                    aria_hidden="true",
                ),
                Div(
                    Span(
                        "Specializations & Foundations",
                        cls="bg-base-100 px-6 text-sm text-base-content/70 uppercase tracking-widest font-semibold select-none",
                    ),
                    cls="relative flex justify-center",
                ),
                cls="relative py-8 max-w-6xl mx-auto",
            ),
            # Specializations Entry Grid
            Div(
                *spec_entries,
                cls="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 max-w-6xl mx-auto mt-6",
            ),
            cls="container mx-auto",
        ),
        id="certifications",
        cls="py-24 bg-base-100 px-4",
    )
