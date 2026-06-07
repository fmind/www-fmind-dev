from fasthtml.common import *

from app.data.projects import OPEN_SOURCE, YOUTUBE_SERIES
from app.icons import icon


def projects_section():
    """Renders the portfolio Projects section."""
    # Build Open Source Projects Grid
    os_project_cards = [
        Div(
            # Top Gradient Accent line
            Div(cls="h-0.5 bg-gradient-to-r from-primary to-blue-400"),
            Div(
                Div(
                    H3(
                        icon(
                            "github" if project.type == "github" else "youtube",
                            cls=f"w-5 h-5 group-hover:scale-110 transition-transform duration-300 {'text-red-600' if project.type == 'youtube' else ''}",
                        ),
                        A(
                            project.title,
                            href=project.href,
                            target="_blank",
                            rel="noopener noreferrer",
                            cls="text-base-content group-hover:text-primary transition-colors hover:underline decoration-2 underline-offset-4 select-text",
                        ),
                        icon(
                            "external_link",
                            cls="w-4 h-4 opacity-0 group-hover:opacity-100 transition-opacity shrink-0",
                        ),
                        cls="card-title text-lg font-bold flex items-center gap-2 select-none",
                    ),
                    P(
                        project.description,
                        cls="text-sm text-base-content/85 mt-4 leading-relaxed text-justify select-text",
                    ),
                ),
                Div(
                    A(
                        "View Code",
                        Span("→", cls="group-hover:translate-x-1 transition-transform inline-block"),
                        href=project.repo if project.repo else project.href,
                        target="_blank",
                        rel="noopener noreferrer",
                        cls="btn btn-sm btn-ghost hover:bg-base-300 font-semibold gap-1",
                    ),
                    cls="card-actions justify-end mt-6",
                ),
                cls="card-body p-6 flex flex-col justify-between",
            ),
            cls="card bg-base-200 shadow-xl hover:shadow-2xl transition-all duration-300 hover:-translate-y-2 border border-base-300 hover:border-primary/30 group overflow-hidden",
        )
        for project in OPEN_SOURCE
    ]

    # Build YouTube Playlists Grid
    yt_playlist_cards = [
        A(
            # Top Gradient Accent line
            Div(cls="h-0.5 bg-gradient-to-r from-primary to-blue-400"),
            Div(
                H3(
                    Span(playlist.title, cls="font-bold text-lg"),
                    icon("external_link", cls="w-4 h-4 opacity-0 group-hover:opacity-100 transition-opacity shrink-0"),
                    cls="card-title text-base-content group-hover:text-primary transition-colors flex items-center gap-1.5 select-none",
                ),
                P(
                    playlist.description,
                    cls="text-sm text-base-content/85 mt-3 leading-relaxed text-justify select-text",
                ),
                Div(
                    Div(
                        playlist.cta,
                        Span("→", cls="group-hover:translate-x-1 transition-transform inline-block"),
                        cls="btn btn-sm btn-ghost hover:bg-base-300 font-semibold gap-1",
                    ),
                    cls="card-actions justify-end mt-6",
                ),
                cls="card-body p-6",
            ),
            href=playlist.url,
            target="_blank",
            rel="noopener noreferrer",
            cls="card bg-base-200 shadow-xl hover:shadow-2xl transition-all duration-300 hover:-translate-y-2 border border-base-300 block group overflow-hidden",
        )
        for playlist in YOUTUBE_SERIES
    ]

    return Section(
        Div(
            # Section Header
            Div(
                H2("Projects", cls="text-4xl md:text-5xl font-heading font-bold mb-4 text-primary"),
                P(
                    "Open-source repositories and educational deep dives created to upskill engineers worldwide.",
                    cls="text-xl text-base-content/70 max-w-3xl mx-auto leading-relaxed text-balance",
                ),
                cls="text-center mb-16",
            ),
            # Open Source Projects Section
            Div(
                H3(
                    A(
                        icon("github", cls="w-6 h-6 text-primary group-hover:text-accent transition-colors"),
                        Span("Open Source Projects"),
                        icon("external_link", cls="w-5 h-5 text-accent shrink-0"),
                        href="https://github.com/fmind",
                        target="_blank",
                        rel="noopener noreferrer",
                        cls="inline-flex items-center gap-2.5 group hover:text-accent transition-colors",
                    ),
                    cls="text-2xl font-bold mb-6 text-base-content",
                ),
                Div(
                    *os_project_cards,
                    cls="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8",
                ),
                cls="mb-16 max-w-6xl mx-auto",
            ),
            # YouTube Series Section
            Div(
                H3(
                    A(
                        icon("youtube", cls="w-6 h-6 text-red-500 group-hover:text-accent transition-colors"),
                        Span("YouTube Series"),
                        icon("external_link", cls="w-5 h-5 text-accent shrink-0"),
                        href="https://www.youtube.com/@BleedingAgent",
                        target="_blank",
                        rel="noopener noreferrer",
                        cls="inline-flex items-center gap-2.5 group hover:text-accent transition-colors",
                    ),
                    cls="text-2xl font-bold mb-6 text-base-content",
                ),
                Div(
                    *yt_playlist_cards,
                    cls="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8",
                ),
                cls="max-w-6xl mx-auto",
            ),
            cls="container mx-auto",
        ),
        id="projects",
        cls="py-24 bg-base-100 px-4",
    )
