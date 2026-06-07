from app.models import CuratedPost, ResearchPaper

PAPERS = [
    ResearchPaper(
        title="Euphony: Harmonious Unification of Cacophonous Anti-Virus Vendor Labels",
        url="https://orbilu.uni.lu/handle/10993/31441",
        venue="MSR 2017 • Mining Software Repositories",
        code="https://github.com/fmind/euphony",
        code_label="Euphony — Label Unification",
    ),
    ResearchPaper(
        title="On the Lack of Consensus in Anti-Virus Decisions",
        url="https://orbilu.uni.lu/handle/10993/27845",
        venue="DIMVA 2016 • Detection of Intrusions and Malware",
        code="https://github.com/fmind/stase",
        code_label="STASE — Statistical Metrics",
    ),
]

POSTS = [
    CuratedPost(
        title="Architecting the AI Agent Platform: A Definitive Guide",
        url="https://fmind.medium.com/architecting-the-ai-agent-platform-a-definitive-guide-405750a3de44",
    ),
    CuratedPost(
        title="Powering Up Your Agent in Production with ADK, OAuth, and Gemini Enterprise",
        url="https://fmind.medium.com/powering-up-your-agent-in-production-with-adk-oauth-and-gemini-enterprise-a52b0716fcba",
    ),
    CuratedPost(
        title="DA2A: The Future of Data Platforms Is Agentic, Distributed and Collaborative",
        url="https://fmind.medium.com/da2a-the-future-of-data-platforms-is-agentic-distributed-and-collaborative-741d5aa96fc4",
    ),
    CuratedPost(
        title="CAG vs RAG: Choosing the Right Strategy for Your AI Application",
        url="https://fmind.medium.com/cag-vs-rag-choosing-the-right-strategy-for-your-ai-application-68dcae85d028",
    ),
    CuratedPost(
        title="Stop Building Rigid AI/ML Pipelines: Embrace Reusable Components",
        url="https://fmind.medium.com/stop-building-rigid-ai-ml-pipelines-embrace-reusable-components-for-flexible-mlops-6e165d837110",
    ),
    CuratedPost(
        title="Poetry Was Good, uv Is Better: An MLOps Migration Story",
        url="https://fmind.medium.com/poetry-was-good-uv-is-better-an-mlops-migration-story-f52bf0c6c703",
    ),
    CuratedPost(
        title="Make Your MLOps Code Base SOLID with Pydantic and Python's ABC",
        url="https://fmind.medium.com/make-your-mlops-code-base-solid-with-pydantic-and-pythons-abc-aeedfe9c3e65",
    ),
    CuratedPost(
        title="How to Configure VS Code for AI/ML and MLOps Development in Python 🐍",
        url="https://fmind.medium.com/how-to-configure-vs-code-for-ai-ml-and-mlops-development-in-python-%EF%B8%8F%EF%B8%8F-8582d8c6ea54",
    ),
]
