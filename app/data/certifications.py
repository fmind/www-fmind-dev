from app.models import CertificationBadge, CertificationEntry

BADGES = [
    CertificationBadge(
        url="https://www.credly.com/badges/8a633d9e-f873-441d-afac-2a8c4d1363b0",
        logo="google.webp",
        title="Professional Cloud Architect",
        issuer="Google Cloud",
        cert_id="f86b41ffd74d4b50947e69876a9274af",
        active=True,
    ),
    CertificationBadge(
        url="https://www.credly.com/badges/6d8416f5-c128-4e25-8557-8ac2289753dd/",
        logo="google.webp",
        title="Professional ML Engineer",
        issuer="Google Cloud",
        cert_id="84c1e188f1054fb8bd53af78060df688",
        active=False,
    ),
    CertificationBadge(
        url="https://credentials.databricks.com/787af897-9bfd-40a2-af97-554bf5a52b74#gs.hfzizp",
        logo="databricks.webp",
        title="Machine Learning Associate",
        issuer="Databricks",
        cert_id="61461287",
        active=False,
    ),
    CertificationBadge(
        url="https://www.credly.com/badges/90e76f14-61b9-45eb-a9d6-0b95c424c6dd",
        logo="microsoft.webp",
        title="Data Scientist Associate",
        issuer="Microsoft Azure",
        cert_id="992564946",
        active=False,
    ),
    CertificationBadge(
        url="https://graphacademy.neo4j.com/c/cd64002c-63a3-4968-87dd-d6cca204a5cd/",
        logo="neo4j.webp",
        title="Graph Data Science",
        issuer="Neo4j",
        cert_id="cd64002c-63a3-4968-87dd-d6cca204a5cd",
        active=False,
    ),
]

SPECIALIZATIONS = [
    CertificationEntry(
        url="https://www.coursera.org/account/accomplishments/specialization/certificate/WLU4DBPSQ4B5",
        logo="google.webp",
        title="Architecting with GKE",
        issuer_details="Google • WLU4DBPSQ4B5",
    ),
    CertificationEntry(
        url="https://www.coursera.org/account/accomplishments/specialization/certificate/EPZ3WQFC423E",
        logo="google.webp",
        title="Cloud Data Engineering",
        issuer_details="Google • EPZ3WQFC423E",
    ),
    CertificationEntry(
        url="https://www.coursera.org/account/accomplishments/specialization/certificate/YSNPABSMV6JL",
        logo="google.webp",
        title="Machine Learning for Trading",
        issuer_details="Google • YSNPABSMV6JL",
    ),
    CertificationEntry(
        url="https://www.coursera.org/account/accomplishments/specialization/certificate/V492QQ4JJKEB",
        logo="google.webp",
        title="Advanced Machine Learning",
        issuer_details="Google • V492QQ4JJKEB",
    ),
    CertificationEntry(
        url="https://www.udemy.com/certificate/UC-XALJEh7G/",
        logo="udemy.webp",
        title="AI: Reinforcement Learning",
        issuer_details="Udemy • UC-XALJEH7G",
    ),
    CertificationEntry(
        url="https://www.udemy.com/certificate/UC-5FM0CC9S/",
        logo="udemy.webp",
        title="Advanced AI: Deep RL",
        issuer_details="Udemy • UC-5FM0CC9S",
    ),
    CertificationEntry(
        url="https://www.coursera.org/account/accomplishments/specialization/certificate/LQ4GHWJ6URBS",
        logo="deeplearning.webp",
        title="TensorFlow Developer",
        issuer_details="DeepLearning • LQ4GHWJ6URBS",
    ),
    CertificationEntry(
        url="https://www.udacity.com/certificate/PV7A7EAA",
        logo="udacity.webp",
        title="Artificial Intelligence",
        issuer_details="Google (Udacity) • PV7A7EAA",
    ),
    CertificationEntry(
        url="https://www.coursera.org/learn/machine-learning",
        logo="stanford.webp",
        title="Machine Learning",
        issuer_details="Stanford (Coursera)",
    ),
]
