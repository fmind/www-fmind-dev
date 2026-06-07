from fasthtml.common import *

from app.icons import icon


def hero_section():
    """Renders the portfolio Hero section."""
    return Div(
        Div(
            Div(
                # Avatar
                Div(
                    Div(
                        Img(
                            src="/static/img/avatar.webp",
                            alt="Médéric Hurier",
                            cls="object-cover w-full h-full",
                            width="192",
                            height="192",
                        ),
                        cls="w-40 h-40 md:w-48 md:h-48 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2 shadow-2xl overflow-hidden mx-auto",
                    ),
                    cls="avatar mb-8",
                ),
                # Headline
                H1(
                    "Médéric Hurier",
                    cls="text-5xl md:text-7xl font-heading font-black mb-4 bg-clip-text text-transparent bg-gradient-to-r from-primary to-blue-400 animate-gradient-x select-text",
                ),
                H2(
                    "Freelancer • AI/ML Architect & Engineer • AI Agents & MLOps •",
                    Br(cls="hidden md:block"),
                    " GCP Professional Cloud Architect • PhD in AI & Computer Security",
                    cls="text-xl md:text-3xl font-heading font-medium text-base-content/90 mb-10 max-w-4xl mx-auto leading-relaxed text-balance select-text",
                ),
                # Actions
                Div(
                    A(
                        "Contact Me",
                        href="mailto:contact@fmind.dev",
                        target="_blank",
                        cls="btn btn-lg bg-primary/5 hover:bg-primary/10 text-primary border border-primary/20 rounded-xl shadow-lg hover:scale-105 active:scale-95 transition-transform hover:shadow-xl hover:shadow-primary/10 px-6 py-3 font-semibold",
                    ),
                    A(
                        "Book a Session",
                        href="https://calendar.google.com/calendar/u/0/appointments/schedules/AcZssZ2ye3X9589PA2xmbV73Iz5J_NbFig6nN651vn6UuYAC-Cs5vBxnQ2L5db9UnAeXmUBQSW1MOobd",
                        target="_blank",
                        rel="noopener noreferrer",
                        cls="btn btn-primary btn-lg bg-primary text-white border-none rounded-xl shadow-lg ring-2 ring-primary/40 hover:scale-105 active:scale-95 transition-transform hover:shadow-xl hover:shadow-primary/20 px-6 py-3 font-semibold",
                    ),
                    cls="flex flex-wrap justify-center gap-4 md:gap-8",
                ),
                # Scroll Down
                Div(
                    A(
                        icon("arrow_down", cls="h-10 w-10"),
                        href="#about",
                        cls="text-primary hover:text-primary-focus transition-colors duration-300",
                        aria_label="Scroll to About section",
                    ),
                    cls="mt-16 animate-bounce flex justify-center",
                ),
                cls="max-w-5xl",
            ),
            cls="hero-content text-center z-10 relative px-4",
        ),
        cls="hero min-h-screen bg-base-200 relative overflow-hidden flex items-center justify-center",
    )
