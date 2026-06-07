from fasthtml.common import *

from app.data.experience import EXPERIENCES


def work_experience_section():
    """Renders the portfolio Work Experience section."""
    experience_items = []
    for item in EXPERIENCES:
        # Build tag badges
        tag_badges = [
            Div(tag, cls="badge badge-outline text-xs border-base-content/25 py-2 px-2.5 font-medium select-none")
            for tag in item.tags
        ]

        experience_items.append(
            Li(
                Div(
                    # Logo Figure Container (Left side)
                    Figure(
                        Div(
                            Img(
                                src=f"/static/img/companies/{item.logo}",
                                alt=item.company,
                                cls="object-contain w-full h-full max-h-20 filter",
                                width="150",
                                height="80",
                            ),
                            cls="rounded-lg p-2 flex items-center justify-center w-full h-full max-h-24 max-w-[150px]",
                        ),
                        cls="w-full sm:w-1/3 p-6 shrink-0 flex items-center justify-center min-h-[150px] bg-slate-50 border-r border-base-200",
                    ),
                    # Experience Info (Right side)
                    Div(
                        Div(
                            H3(
                                item.title,
                                cls="card-title text-lg md:text-xl font-bold leading-tight text-base-content select-text",
                            ),
                            P(item.company, cls="text-sm font-semibold text-primary mt-1 select-text"),
                            P(
                                item.description,
                                cls="text-base-content/80 text-sm mt-4 text-justify leading-relaxed select-text",
                            ),
                        ),
                        # Tag Badges
                        Div(
                            *tag_badges,
                            cls="card-actions justify-start gap-1.5 mt-6",
                        ),
                        cls="card-body w-full sm:w-2/3 p-6 flex flex-col justify-between",
                    ),
                    style=f"border-left-color: {item.hex};",
                    cls="card card-side bg-base-100 shadow-xl border-l-4 hover:shadow-2xl transition-all duration-300 hover:-translate-y-2 h-full flex flex-col sm:flex-row overflow-hidden group",
                ),
                cls="h-full",
            )
        )

    return Section(
        Div(
            # Section Header
            Div(
                H2(
                    "Work ",
                    Span("Experience", cls="text-primary"),
                    cls="text-4xl md:text-5xl font-heading font-bold mb-4 text-primary",
                ),
                P(
                    "A track record of engineering scale, reliability, and security for industry-leading institutions.",
                    cls="text-xl text-base-content/70 max-w-3xl mx-auto leading-relaxed text-balance",
                ),
                cls="text-center mb-16",
            ),
            # Experience Timeline Grid
            Ul(
                *experience_items,
                cls="grid grid-cols-1 md:grid-cols-2 gap-8 max-w-6xl mx-auto",
            ),
            cls="container mx-auto",
        ),
        id="work-experience",
        cls="py-24 bg-base-200 px-4",
    )
