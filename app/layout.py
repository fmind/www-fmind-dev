import hashlib
import json
from datetime import datetime
from pathlib import Path

from fasthtml.common import *

from app.icons import icon


# Global CSS version computed at startup for cache busting
def _get_css_version() -> str:
    css_path = Path("static/dist/styles.css")
    if css_path.exists():
        return hashlib.sha256(css_path.read_bytes()).hexdigest()[:8]
    return "1.0.0"


CSS_VERSION = _get_css_version()


def get_headers() -> list:
    """Returns the list of global headers matching original index.html metadata and preloads."""
    theme_js = """
      (function () {
        function setTheme(isDark) {
          localStorage.theme = isDark ? 'dark' : 'light';
          document.documentElement.classList.toggle('dark', isDark);
          document.documentElement.setAttribute('data-theme', isDark ? 'dark' : 'light');
          document.querySelectorAll('meta[name="theme-color"]').forEach((meta) => {
            meta.setAttribute('content', isDark ? '#1d232a' : '#ffffff');
          });
        }
        var isDark =
          localStorage.theme === 'dark' ||
          (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches);
        setTheme(isDark);
        window.setTheme = setTheme;
      })();
    """

    person_schema = {
        "@context": "https://schema.org",
        "@type": "Person",
        "name": "Médéric Hurier",
        "alternateName": "Fmind",
        "url": "https://www.fmind.dev",
        "image": "https://www.fmind.dev/static/img/avatar.webp",
        "email": "mailto:contact@fmind.dev",
        "jobTitle": "AI/ML Architect & Engineer",
        "description": (
            "Freelance AI/ML Architect & Engineer specializing in AI Agents, Generative AI, MLOps, and Google Cloud "
            "Platform. PhD in AI & Computer Security."
        ),
        "nationality": {"@type": "Country", "name": "France"},
        "knowsLanguage": ["fr", "en"],
        "knowsAbout": [
            "Artificial Intelligence",
            "Machine Learning",
            "MLOps",
            "AI Agents",
            "Generative AI",
            "Computer Security",
            "Google Cloud Platform",
        ],
        "alumniOf": {
            "@type": "CollegeOrUniversity",
            "name": "University of Luxembourg",
            "sameAs": "https://wwwen.uni.lu/snt/people/mederic_hurier",
        },
        "hasCredential": [
            {
                "@type": "EducationalOccupationalCredential",
                "name": "Professional Cloud Architect",
                "credentialCategory": "certification",
                "recognizedBy": {"@type": "Organization", "name": "Google Cloud"},
            }
        ],
        "sameAs": [
            "https://www.linkedin.com/in/fmind-dev/",
            "https://x.com/fmind_dev",
            "https://github.com/fmind",
            "https://fmind.medium.com/",
            "https://www.youtube.com/@fmind-dev",
        ],
    }

    website_schema = {
        "@context": "https://schema.org",
        "@type": "WebSite",
        "name": "www.fmind.dev",
        "url": "https://www.fmind.dev",
        "description": (
            "Freelance AI/ML Architect & Engineer specializing in AI Agents, Generative AI, MLOps, and Google Cloud "
            "Platform. PhD in AI & Computer Security."
        ),
        "author": {"@type": "Person", "name": "Médéric Hurier"},
    }

    return [
        Title("www.fmind.dev — Médéric Hurier | AI/ML Architect & Engineer"),
        Meta(charset="utf-8"),
        Meta(name="viewport", content="width=device-width, initial-scale=1"),
        Meta(name="theme-color", content="#ffffff"),
        Meta(name="author", content="Médéric Hurier (Fmind)"),
        Meta(
            name="description",
            content=(
                "Freelance AI/ML Architect & Engineer specializing in AI Agents, Generative AI, MLOps, and Google Cloud "
                "Platform. PhD in AI & Computer Security."
            ),
        ),
        Meta(
            name="keywords",
            content=(
                "AI, Machine Learning, MLOps, Artificial Intelligence, AI Agents, Agentic AI, Generative AI, Google Cloud, "
                "GCP, Python, Freelance, Luxembourg"
            ),
        ),
        Meta(name="robots", content="index, follow"),
        Link(rel="me", href="https://www.linkedin.com/in/fmind-dev/"),
        Link(rel="me", href="https://x.com/fmind_dev"),
        Link(rel="me", href="https://fmind.medium.com"),
        Link(rel="me", href="https://github.com/fmind"),
        Link(rel="me", href="https://www.youtube.com/@fmind-dev"),
        Meta(property="og:type", content="website"),
        Meta(property="og:url", content="https://www.fmind.dev/"),
        Meta(property="og:title", content="www.fmind.dev — Médéric Hurier | AI/ML Architect & Engineer"),
        Meta(property="og:site_name", content="www.fmind.dev"),
        Meta(
            property="og:description",
            content=(
                "Freelance AI/ML Architect & Engineer specializing in AI Agents, Generative AI, MLOps, and Google Cloud "
                "Platform. PhD in AI & Computer Security."
            ),
        ),
        Meta(property="og:image", content="https://www.fmind.dev/static/img/og-image.jpg"),
        Meta(property="og:image:width", content="1200"),
        Meta(property="og:image:height", content="630"),
        Meta(property="og:image:alt", content="www.fmind.dev — AI/ML Architect & Engineer Banner"),
        Meta(property="og:locale", content="en_US"),
        Meta(name="twitter:card", content="summary_large_image"),
        Meta(name="twitter:site", content="@fmind_dev"),
        Meta(name="twitter:creator", content="@fmind_dev"),
        Meta(name="twitter:title", content="www.fmind.dev — Médéric Hurier | AI/ML Architect & Engineer"),
        Meta(
            name="twitter:description",
            content=(
                "Freelance AI/ML Architect & Engineer specializing in AI Agents, Generative AI, MLOps, and Google Cloud "
                "Platform. PhD in AI & Computer Security."
            ),
        ),
        Meta(name="twitter:image", content="https://www.fmind.dev/static/img/og-image.jpg"),
        Meta(name="twitter:image:alt", content="www.fmind.dev — AI/ML Architect & Engineer banner"),
        Script(theme_js),
        # Font & image preloads
        Link(
            rel="preload",
            href="/static/fonts/Outfit-Variable.woff2",
            crossorigin="anonymous",
            **{"as": "font", "type": "font/woff2"},
        ),
        Link(
            rel="preload",
            href="/static/fonts/Inter-Variable.woff2",
            crossorigin="anonymous",
            **{"as": "font", "type": "font/woff2"},
        ),
        Link(rel="preload", href="/static/img/avatar.webp", fetchpriority="high", **{"as": "image"}),
        Link(rel="icon", type="image/x-icon", href="/static/favicon.ico"),
        Link(rel="apple-touch-icon", sizes="180x180", href="/static/img/favicons/apple-touch-icon.png"),
        Link(rel="icon", type="image/png", sizes="32x32", href="/static/img/favicons/favicon-32x32.png"),
        Link(rel="icon", type="image/png", sizes="16x16", href="/static/img/favicons/favicon-16x16.png"),
        Link(rel="manifest", href="/static/site.webmanifest"),
        # Local CSS & JS libraries (External cached stylesheet with cache busting)
        Link(rel="stylesheet", href=f"/static/dist/styles.css?v={CSS_VERSION}"),
        Script(src="/static/js/htmx.min.js", defer=True),
        Script(src="/static/js/alpine.min.js", defer=True),
        Script(json.dumps(person_schema), type="application/ld+json"),
        Script(json.dumps(website_schema), type="application/ld+json"),
    ]


def navigation_bar():
    """Renders the responsive header navigation bar using Alpine.js and DaisyUI."""
    # Links matching the Angular header
    nav_links = [
        {"href": "#about", "label": "About"},
        {"href": "#work-experience", "label": "Work Experience"},
        {"href": "#certifications", "label": "Certifications"},
        {"href": "#publications", "label": "Publications"},
        {"href": "#projects", "label": "Projects"},
        {"href": "#services", "label": "Services"},
    ]

    social_links = [
        {"href": "https://www.linkedin.com/in/fmind-dev/", "label": "LinkedIn", "icon": "linkedin"},
        {"href": "https://x.com/fmind_dev", "label": "X (Twitter)", "icon": "x"},
        {"href": "https://fmind.medium.com/", "label": "Medium", "icon": "medium"},
        {"href": "https://github.com/fmind", "label": "GitHub", "icon": "github"},
        {"href": "https://www.youtube.com/@fmind-dev", "label": "YouTube", "icon": "youtube"},
    ]

    # Mobile menu options
    mobile_dropdown_items = [
        Li(A(link["label"], href=link["href"], cls="font-medium hover:text-primary py-2", **{"@click": "open = false"}))
        for link in nav_links
    ]

    # Desktop menu options
    desktop_menu_items = [
        Li(
            A(
                link["label"],
                href=link["href"],
                cls=(
                    "relative font-medium hover:text-primary transition-colors duration-300 "
                    "after:content-[''] after:absolute after:w-full after:scale-x-0 after:h-0.5 after:bottom-0 after:left-0 "
                    "after:bg-primary after:origin-bottom-right after:transition-transform after:duration-300 "
                    "hover:after:scale-x-100 hover:after:origin-bottom-left py-2 px-3"
                ),
            )
        )
        for link in nav_links
    ]

    # Social icons
    social_icons = [
        A(
            icon(social["icon"], cls="h-5 w-5 md:h-6 md:w-6"),
            href=social["href"],
            target="_blank",
            rel="noopener noreferrer",
            aria_label=social["label"],
            cls="btn btn-ghost btn-circle btn-sm md:btn-md hover:bg-transparent hover:text-accent transition-colors",
        )
        for social in social_links
    ]

    return Nav(
        Div(
            # Start section (Mobile menu toggle + brand)
            Div(
                # Alpine wrapper for mobile menu state
                Div(
                    Button(
                        icon("menu", cls="h-5 w-5"),
                        cls="btn btn-ghost btn-circle",
                        aria_label="Menu",
                        **{"@click": "open = !open"},
                    ),
                    Ul(
                        *mobile_dropdown_items,
                        cls="menu menu-sm absolute left-0 mt-3 z-[50] p-2 shadow bg-base-100 rounded-box w-52 border border-base-200",
                        **{"x-show": "open", "@click.outside": "open = false"},
                    ),
                    cls="lg:hidden relative",
                    **{"x-data": "{ open: false }"},
                ),
                A(
                    Img(
                        src="/static/img/logo.webp",
                        alt="www.fmind.dev logo",
                        cls="w-10 h-10 md:w-12 md:h-12 rounded-full object-cover",
                        width="48",
                        height="48",
                    ),
                    Span("www.fmind.dev", cls="hidden sm:inline"),
                    href="#",
                    cls="btn btn-ghost text-xl md:text-2xl font-bold font-heading text-primary normal-case gap-3 sm:gap-4 px-2",
                ),
                Ul(*desktop_menu_items, cls="menu menu-horizontal px-1 hidden lg:flex gap-1"),
                cls="flex-1 flex items-center gap-4",
            ),
            # End section (Social icons + theme toggle)
            Div(
                Div(*social_icons, cls="flex gap-1 md:gap-3"),
                Div(cls="border-l border-base-content/20 h-6 md:h-8 mx-1 md:mx-2"),
                # Theme toggle button utilizing Alpine reactive state
                Button(
                    # Moon icon (shows in light theme)
                    Span(icon("moon", cls="h-5 w-5 md:h-6 md:w-6"), **{"x-show": "theme === 'light'"}),
                    # Sun icon (shows in dark theme)
                    Span(icon("sun", cls="h-5 w-5 md:h-6 md:w-6"), **{"x-show": "theme === 'dark'"}),
                    cls="btn btn-ghost btn-circle btn-sm md:btn-md bg-base-200/60 border border-base-content/10",
                    aria_label="Toggle Dark Mode",
                    **{"@click": "toggleTheme()"},
                ),
                cls="flex-none flex items-center gap-2 md:gap-4",
            ),
            cls="container mx-auto flex px-4",
        ),
        cls="navbar bg-base-100/90 backdrop-blur-md border-b border-base-200",
    )


def page_footer():
    """Renders the standard footer component."""
    return Footer(
        Aside(
            P(
                "Industrialize Intelligence with Craft and Passion.",
                cls="font-heading font-semibold text-lg md:text-xl italic text-slate-100",
            ),
            P(f"© {datetime.now().year} Médéric Hurier (Fmind). All rights reserved.", cls="text-sm"),
            cls="space-y-2",
        ),
        cls="footer footer-center p-8 bg-slate-950 text-slate-400 border-t border-slate-900",
    )


def page_layout(content: Any):
    """Wraps body content in the main header and footer shell, with Alpine.js theme initialization."""
    return Div(
        # Header / Navbar
        Header(navigation_bar(), cls="sticky top-0 z-40"),
        # Main content area
        Main(content, id="main-content", cls="min-h-screen"),
        # Footer
        page_footer(),
        cls="min-h-screen bg-base-100 text-base-content font-body antialiased transition-colors duration-300",
        **{
            "x-data": (
                "{ theme: localStorage.theme || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'), "
                "toggleTheme() { this.theme = this.theme === 'light' ? 'dark' : 'light'; localStorage.theme = this.theme; "
                "document.documentElement.classList.toggle('dark', this.theme === 'dark'); "
                "document.documentElement.setAttribute('data-theme', this.theme); "
                "document.querySelectorAll('meta[name=\\'theme-color\\']').forEach(m => m.setAttribute('content', this.theme === 'dark' ? '#1d232a' : '#ffffff')); } }"
            )
        },
    )
