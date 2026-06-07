from fasthtml.common import *

from app.data.about import BIOGRAPHY, EXPERTISE


def about_section():
    """Renders the portfolio About section."""
    # Build Biography Paragraphs (using NotStr to preserve HTML formatting)
    bio_paras = [P(NotStr(para), cls="select-text") for para in BIOGRAPHY]

    # Build Core Expertise Cards
    expertise_cards = [
        Div(
            # Top Gradient Bar
            Div(cls=f"h-1 bg-gradient-to-r {card.gradient}"),
            Div(
                Div(card.emoji, cls="text-5xl mb-3 group-hover:scale-110 transition-transform duration-300"),
                H4(card.title, cls="font-bold text-lg text-base-content"),
                P(card.description, cls="text-sm text-base-content/70 mt-2 leading-relaxed"),
                cls="card-body p-6 items-center text-center",
            ),
            cls="card bg-base-200/50 backdrop-blur-sm hover:bg-base-200 transition-all duration-300 hover:-translate-y-2 border border-base-300 hover:border-primary/30 overflow-hidden h-full shadow-xl hover:shadow-2xl group",
        )
        for card in EXPERTISE
    ]

    return Section(
        Div(
            # Section Header
            Div(
                H2(
                    "About ",
                    Span("Me", cls="text-primary"),
                    cls="text-4xl md:text-5xl font-heading font-bold mb-4 text-primary",
                ),
                P(
                    "Transforming AI/ML initiatives into secure, scalable, and high-impact solutions.",
                    cls="text-xl text-base-content/70 max-w-3xl mx-auto leading-relaxed text-balance",
                ),
                cls="text-center mb-12",
            ),
            # Biography
            Div(
                Div(
                    Div(
                        *bio_paras,
                        cls="prose prose-lg text-base-content/80 text-justify mx-auto leading-loose [&_p]:leading-loose space-y-6 select-text",
                    ),
                    cls="md:col-span-10 md:col-start-2 space-y-6",
                ),
                cls="grid md:grid-cols-12 gap-8 mb-16",
            ),
            # Core Expertise Divider
            H3(
                Span(cls="w-2 h-8 bg-primary rounded-full"),
                "Core Expertise & Skills",
                cls="text-xl font-bold mb-10 flex items-center justify-center gap-2",
            ),
            # Cards Grid
            Div(
                *expertise_cards,
                cls="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 max-w-6xl mx-auto",
            ),
            cls="container mx-auto",
        ),
        id="about",
        cls="py-24 bg-base-100 overflow-hidden px-4",
    )
