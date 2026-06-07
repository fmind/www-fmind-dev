from fasthtml.common import *

from app.data.publications import PAPERS, POSTS
from app.icons import icon


def publications_section():
    """Renders the portfolio Publications section."""
    # Build Academic Papers List
    academic_paper_cards = []
    for paper in PAPERS:
        # Code link button if present
        code_action = ""
        if paper.code:
            code_action = Div(
                A(
                    icon("github", cls="w-4 h-4 mr-1.5"),
                    paper.code_label,
                    href=paper.code,
                    target="_blank",
                    rel="noopener noreferrer",
                    cls="btn btn-ghost btn-xs sm:btn-sm border border-base-300 hover:bg-base-200",
                ),
                cls="card-actions shrink-0",
            )

        academic_paper_cards.append(
            Div(
                Div(
                    A(
                        paper.title,
                        icon(
                            "external_link",
                            cls="w-4 h-4 text-accent opacity-0 group-hover:opacity-100 transition-opacity shrink-0",
                        ),
                        href=paper.url,
                        target="_blank",
                        rel="noopener noreferrer",
                        cls="font-bold text-base md:text-md group-hover:text-accent transition-colors inline-flex items-center gap-1.5 select-text",
                    ),
                    Div(
                        P(paper.venue, cls="text-base-content/60 text-xs md:text-sm truncate min-w-0 select-text"),
                        code_action,
                        cls="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mt-3",
                    ),
                    cls="card-body p-5",
                ),
                cls="group card bg-base-100 shadow-sm hover:shadow-md transition-all duration-300 border border-base-300 border-l-2 border-l-transparent hover:border-l-primary/50",
            )
        )

    # Build Curated Posts Links
    curated_post_links = [
        A(
            Div(icon("medium", cls="w-5 h-5"), cls="mt-0.5 shrink-0 text-accent"),
            Span(
                post.title,
                cls="text-sm font-semibold text-base-content/85 group-hover:text-primary transition-colors leading-snug select-text",
            ),
            icon(
                "external_link",
                cls="w-4 h-4 shrink-0 text-accent opacity-0 group-hover:opacity-100 transition-opacity ml-auto",
            ),
            href=post.url,
            target="_blank",
            rel="noopener noreferrer",
            cls="flex items-start gap-3 p-4 hover:bg-base-200 transition-colors duration-200 group",
        )
        for post in POSTS
    ]

    return Section(
        Div(
            # Section Header
            Div(
                H2("Publications", cls="text-4xl md:text-5xl font-heading font-bold mb-4 text-primary"),
                P(
                    "Sharing research, insights, and structural guidelines with the global engineering community.",
                    cls="text-xl text-base-content/70 max-w-3xl mx-auto leading-relaxed text-balance",
                ),
                cls="text-center mb-16",
            ),
            # Content Grid: Left Column (Academic) & Right Column (Featured Articles)
            Div(
                # Left Column: Thesis & Academic Papers
                Div(
                    # PhD Thesis Highlight Card
                    Div(
                        Div(
                            Div("PhD Thesis", cls="badge badge-primary badge-outline font-semibold mb-3 select-none"),
                            A(
                                "Creating better ground truth to further understand Android malware",
                                icon(
                                    "external_link",
                                    cls="w-5 h-5 text-accent opacity-0 group-hover:opacity-100 transition-opacity shrink-0",
                                ),
                                href="https://orbilu.uni.lu/handle/10993/39903",
                                target="_blank",
                                rel="noopener noreferrer",
                                cls="card-title text-xl md:text-2xl font-bold mb-2 group-hover:text-accent transition-colors inline-flex items-center gap-2 select-text",
                            ),
                            P(
                                "University of Luxembourg (SNT) & Google, 2019",
                                cls="text-base-content/70 italic text-sm mb-4 select-text",
                            ),
                            P(
                                "AI/ML models are only as good as the data they learn from — yet Android malware ground truths are "
                                "notoriously unreliable. This thesis tackles the problem by benchmarking antivirus engines, harmonizing "
                                "their conflicting labels, and mining large-scale datasets to characterize malicious behavior.",
                                cls="mb-6 text-justify text-sm md:text-base leading-relaxed text-base-content/85 select-text",
                            ),
                            Div(
                                A(
                                    icon("github", cls="w-4 h-4 mr-1.5"),
                                    "SetValX — Android Malware Processing",
                                    href="https://github.com/fmind/servalx",
                                    target="_blank",
                                    rel="noopener noreferrer",
                                    cls="btn btn-ghost btn-xs sm:btn-sm border border-base-300 hover:bg-base-200",
                                ),
                                A(
                                    icon("github", cls="w-4 h-4 mr-1.5"),
                                    "APKWorkers — Distributed Android Analysis",
                                    href="https://github.com/fmind/apkworkers",
                                    target="_blank",
                                    rel="noopener noreferrer",
                                    cls="btn btn-ghost btn-xs sm:btn-sm border border-base-300 hover:bg-base-200",
                                ),
                                cls="card-actions flex flex-wrap gap-2",
                            ),
                            cls="card-body p-6 md:p-8",
                        ),
                        cls="group card bg-base-100 shadow-xl border-l-4 border-primary hover:shadow-2xl transition-all duration-300 hover:-translate-y-1",
                    ),
                    # Academic Research List
                    Div(
                        H3(
                            icon("academic_cap", cls="w-6 h-6 text-primary"),
                            "Academic Research",
                            cls="text-2xl font-bold mb-6 text-base-content flex items-center gap-2",
                        ),
                        Div(
                            *academic_paper_cards,
                            cls="space-y-4",
                        ),
                    ),
                    cls="lg:col-span-7 space-y-8",
                ),
                # Right Column: Curated Posts
                Div(
                    H3(
                        Span(cls="w-2 h-8 bg-primary rounded-full"),
                        "Featured Posts",
                        cls="text-2xl font-bold text-base-content flex items-center gap-2",
                    ),
                    Div(
                        Div(
                            *curated_post_links,
                            cls="card-body p-2 divide-y divide-base-200",
                        ),
                        cls="card bg-base-100 shadow-xl border border-base-300 overflow-hidden",
                    ),
                    A(
                        icon("medium", cls="w-5 h-5 mr-2"),
                        "More publications on Medium",
                        icon("external_link", cls="w-4 h-4 ml-1.5"),
                        href="https://fmind.medium.com/",
                        target="_blank",
                        rel="noopener noreferrer",
                        cls="btn btn-outline btn-primary btn-md w-full shadow-md rounded-xl hover:scale-[1.01] transition-transform font-bold",
                    ),
                    cls="lg:col-span-5 space-y-6",
                ),
                cls="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start max-w-6xl mx-auto",
            ),
            cls="container mx-auto",
        ),
        id="publications",
        cls="py-24 bg-base-200 px-4",
    )
