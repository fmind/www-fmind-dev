from app.models import ExpertiseCard

BIOGRAPHY = [
    (
        "As a freelance <strong>AI/ML Architect &amp; Engineer</strong> with a Ph.D. in AI and Computer Security, "
        "I specialize in designing autonomous agentic workflows and production-ready MLOps environments. My focus "
        "is on driving real-world business outcomes by ensuring reliability and observability in complex systems."
    ),
    (
        "Certified as a <strong>Google Cloud Professional Cloud Architect</strong>, I combine deep academic "
        "research with a pragmatic, cartesian approach to engineering. I have helped institutions like Google, "
        "Decathlon, BNP Paribas, ArcelorMittal, and the European Commission grow their AI/ML projects to "
        "enterprise-grade solutions in production."
    ),
]

EXPERTISE = [
    ExpertiseCard(
        title="Agentic Orchestration",
        emoji="🤖",
        gradient="from-accent to-blue-500",
        description="Building autonomous systems, reasoning engines, and reliable agentic workflows.",
    ),
    ExpertiseCard(
        title="Production MLOps",
        emoji="🚀",
        gradient="from-blue-500 to-cyan-500",
        description="Robust deployment patterns on GCP, AWS, Azure, and Databricks for reliability.",
    ),
    ExpertiseCard(
        title="Security-First AI",
        emoji="🛡️",
        gradient="from-rose-500 to-red-500",
        description="Leveraging Ph.D. background to build secure and trustworthy AI systems.",
    ),
    ExpertiseCard(
        title="Technical Strategy",
        emoji="🧭",
        gradient="from-orange-500 to-amber-500",
        description="Translating complex AI capabilities into clear, scalable architectural roadmaps.",
    ),
    ExpertiseCard(
        title="Data Science & ML",
        emoji="📊",
        gradient="from-cyan-500 to-blue-500",
        description="Engineering machine learning models and data solutions for enterprise scale.",
    ),
    ExpertiseCard(
        title="Python Development",
        emoji="🐍",
        gradient="from-emerald-500 to-teal-500",
        description="Building scalable applications, libraries, and tools with modern standards.",
    ),
]
