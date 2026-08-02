package templates

// Portfolio Data
var METADATA = Metadata{
	Name:              "Médéric Hurier",
	AlternateName:     "Fmind",
	SiteName:          "Fmind",
	Title:             "Médéric Hurier (Fmind) | AI Architect (PhD) • Freelancer",
	JobTitle:          "AI Architect (PhD) • Freelancer",
	HeadlinePrimary:   "AI Architect (PhD) • VC Expert Advisor • AAIF Ambassador • GCP Certified Cloud Architect • AI Agents, MLOps & Security",
	HeadlineSecondary: "",
	Description:       "Freelance AI Architect (PhD), VC Expert Advisor, AAIF Ambassador, and GCP Certified Cloud Architect specializing in production AI Agents, MLOps, and security.",
	Keywords: []string{
		"AI",
		"Machine Learning",
		"MLOps",
		"Artificial Intelligence",
		"AI Agents",
		"Agentic AI",
		"Generative AI",
		"Agentic AI Foundation",
		"AAIF",
		"Google Cloud",
		"GCP",
		"AI Security",
		"Venture Advisory",
		"Python",
		"Freelance",
		"Luxembourg",
	},
	Email:         "contact@fmind.dev",
	CalendarURL:   "https://calendar.google.com/calendar/u/0/appointments/schedules/AcZssZ2ye3X9589PA2xmbV73Iz5J_NbFig6nN651vn6UuYAC-Cs5vBxnQ2L5db9UnAeXmUBQSW1MOobd",
	SiteURL:       "https://www.fmind.dev",
	TwitterHandle: "@fmind_dev",
	Socials: []SocialLink{
		{Name: "LinkedIn", URL: "https://www.linkedin.com/in/fmind-dev/", Icon: "linkedin", Header: true},
		{Name: "X (Twitter)", URL: "https://x.com/fmind_dev", Icon: "x", Header: true},
		{Name: "Bluesky", URL: "https://bsky.app/profile/fmind-dev.bsky.social", Icon: "bluesky"},
		{Name: "Medium", URL: "https://fmind.medium.com/", Icon: "medium", Header: true},
		{Name: "GitHub", URL: "https://github.com/fmind", Icon: "github", Header: true},
		{Name: "YouTube", URL: "https://www.youtube.com/@fmind-dev", Icon: "youtube", Header: true},
	},
}

var BIOGRAPHY = []string{
	"I am a **freelance AI Architect** with a **PhD in AI and Computer Security**. I design and industrialize **AI agents**, **MLOps platforms**, and **secure cloud foundations**, turning fast-moving research into dependable production capabilities with clear controls, observability, and measurable outcomes.",
	"My work spans strategy and delivery: enterprise agent platforms at **Decathlon**, European fraud detection for the **European Commission**, and Android malware research with **Google**. I have also delivered AI, data, and security initiatives for BNP Paribas, ArcelorMittal, SFEIR, Clearstream, and the University of Luxembourg.",
	"Beyond client work, I serve as an **AAIF Ambassador** and Luxembourg organizer and contribute to the **33N Ventures Expert Advisory Board**. As a **Google Cloud Professional Cloud Architect**, I bring a pragmatic, security-first approach to systems that must operate reliably at scale.",
}

var LEADERSHIP = []LeadershipRole{
	{
		Role:         "Agentic AI Foundation Ambassador",
		Organization: "The Linux Foundation",
		Description:  "Selected to serve as an AAIF Ambassador and help grow the open agentic AI community.",
		URL:          "https://www.credly.com/badges/aaf051e1-202f-4b0f-bfc6-a23a3ef2e2a2",
	},
	{
		Role:         "Expert Advisory Board Member",
		Organization: "33N Ventures",
		Description:  "Contributing AI Agents, MLOps, and security expertise to 33N's venture advisory network.",
		URL:          "https://33n.vc/team",
	},
	{
		Role:         "Local Community Organizer",
		Organization: "AAIF Community Luxembourg",
		Description:  "Organizing Luxembourg's local practitioner community and events around agentic AI.",
		URL:          "https://luma.com/aaif-luxembourg",
	},
}

var EXPERTISE = []ExpertiseCard{
	{
		Title:       "Agentic Orchestration",
		Emoji:       "🤖",
		Description: "Building autonomous systems, reasoning engines, and reliable agentic workflows.",
	},
	{
		Title:       "Production MLOps",
		Emoji:       "🚀",
		Description: "Robust deployment patterns on GCP, AWS, Azure, and Databricks for reliability.",
	},
	{
		Title:       "Security-First AI",
		Emoji:       "🛡️",
		Description: "Leveraging a PhD background to build secure and trustworthy AI systems.",
	},
	{
		Title:       "Technical Strategy",
		Emoji:       "🧭",
		Description: "Translating complex AI capabilities into clear, scalable architectural roadmaps.",
	},
	{
		Title:       "Data Science & ML",
		Emoji:       "📊",
		Description: "Engineering machine learning models and data solutions for enterprise scale.",
	},
	{
		Title:       "Python Development",
		Emoji:       "🐍",
		Description: "Building scalable applications, libraries, and tools with modern standards.",
	},
}

var BADGES = []CertificationBadge{
	{
		URL:    "https://www.credly.com/badges/aaf051e1-202f-4b0f-bfc6-a23a3ef2e2a2",
		Logo:   "aaif.webp",
		Title:  "Agentic AI Foundation Ambassador",
		Issuer: "The Linux Foundation",
		CertID: "2026",
		Active: true,
	},
	{
		URL:    "https://www.credly.com/badges/8a633d9e-f873-441d-afac-2a8c4d1363b0",
		Logo:   "google.webp",
		Title:  "Professional Cloud Architect",
		Issuer: "Google Cloud",
		CertID: "f86b41ffd74d4b50947e69876a9274af",
		Active: true,
	},
	{
		URL:    "https://www.credly.com/badges/6d8416f5-c128-4e25-8557-8ac2289753dd/",
		Logo:   "google.webp",
		Title:  "Professional ML Engineer",
		Issuer: "Google Cloud",
		CertID: "84c1e188f1054fb8bd53af78060df688",
		Active: false,
	},
	{
		URL:    "https://credentials.databricks.com/787af897-9bfd-40a2-af97-554bf5a52b74#gs.hfzizp",
		Logo:   "databricks.webp",
		Title:  "Machine Learning Associate",
		Issuer: "Databricks",
		CertID: "61461287",
		Active: false,
	},
	{
		URL:    "https://www.credly.com/badges/90e76f14-61b9-45eb-a9d6-0b95c424c6dd",
		Logo:   "microsoft.webp",
		Title:  "Data Scientist Associate",
		Issuer: "Microsoft Azure",
		CertID: "992564946",
		Active: false,
	},
	{
		URL:    "https://graphacademy.neo4j.com/c/cd64002c-63a3-4968-87dd-d6cca204a5cd/",
		Logo:   "neo4j.webp",
		Title:  "Graph Data Science",
		Issuer: "Neo4j",
		CertID: "cd64002c-63a3-4968-87dd-d6cca204a5cd",
		Active: false,
	},
}

var SPECIALIZATIONS = []CertificationEntry{
	{
		URL:           "https://www.coursera.org/account/accomplishments/specialization/certificate/WLU4DBPSQ4B5",
		Logo:          "google.webp",
		Title:         "Architecting with GKE",
		IssuerDetails: "Google • WLU4DBPSQ4B5",
	},
	{
		URL:           "https://www.coursera.org/account/accomplishments/specialization/certificate/EPZ3WQFC423E",
		Logo:          "google.webp",
		Title:         "Cloud Data Engineering",
		IssuerDetails: "Google • EPZ3WQFC423E",
	},
	{
		URL:           "https://www.coursera.org/account/accomplishments/specialization/certificate/YSNPABSMV6JL",
		Logo:          "google.webp",
		Title:         "Machine Learning for Trading",
		IssuerDetails: "Google • YSNPABSMV6JL",
	},
	{
		URL:           "https://www.coursera.org/account/accomplishments/specialization/certificate/V492QQ4JJKEB",
		Logo:          "google.webp",
		Title:         "Advanced Machine Learning",
		IssuerDetails: "Google • V492QQ4JJKEB",
	},
	{
		URL:           "https://www.udemy.com/certificate/UC-XALJEh7G/",
		Logo:          "udemy.webp",
		Title:         "AI: Reinforcement Learning",
		IssuerDetails: "Udemy • UC-XALJEH7G",
	},
	{
		URL:           "https://www.udemy.com/certificate/UC-5FM0CC9S/",
		Logo:          "udemy.webp",
		Title:         "Advanced AI: Deep RL",
		IssuerDetails: "Udemy • UC-5FM0CC9S",
	},
	{
		URL:           "https://www.coursera.org/account/accomplishments/specialization/certificate/LQ4GHWJ6URBS",
		Logo:          "deeplearning.webp",
		Title:         "TensorFlow Developer",
		IssuerDetails: "DeepLearning • LQ4GHWJ6URBS",
	},
	{
		URL:           "https://www.udacity.com/certificate/PV7A7EAA",
		Logo:          "udacity.webp",
		Title:         "Artificial Intelligence",
		IssuerDetails: "Google (Udacity) • PV7A7EAA",
	},
	{
		URL:           "https://www.coursera.org/learn/machine-learning",
		Logo:          "stanford.webp",
		Title:         "Machine Learning",
		IssuerDetails: "Stanford (Coursera)",
	},
}

var EXPERIENCES = []WorkExperience{
	{
		Company:     "Decathlon",
		Logo:        "decathlon.webp",
		Title:       "AI/ML Architect",
		BrandColor:  "#3643BA",
		Description: "Design and implement enterprise-scale Agents and MLOps platforms for AI/ML industrialization.",
		Tags:        []string{"AI/ML", "Agents", "Gen AI", "MLOps"},
	},
	{
		Company:     "European Commission",
		Logo:        "european-commission.webp",
		Title:       "AI/ML Engineer",
		BrandColor:  "#004494",
		Description: "Contributed to Arachne, the European fraud detection system to ensure financial integrity.",
		Tags:        []string{"AI/ML", "Fraud Detection", "Public Sector"},
	},
	{
		Company:     "ArcelorMittal",
		Logo:        "arcelor-mittal.webp",
		Title:       "Data Scientist",
		BrandColor:  "#F47D30",
		Description: "Trained and optimized AI/ML models for steel price recommendations in industrial markets.",
		Tags:        []string{"AI/ML", "Data", "Forecasting", "Industry"},
	},
	{
		Company:     "BNP Paribas",
		Logo:        "bgl-bnp-paribas.webp",
		Title:       "Project Manager",
		BrandColor:  "#00915E",
		Description: "Supervised the development of advanced transformer models for banking applications.",
		Tags:        []string{"NLP", "Finance", "Project Management"},
	},
	{
		Company:     "Google",
		Logo:        "google.webp",
		Title:       "Research Partner",
		BrandColor:  "#4285F4",
		Description: "Collaborated with the Android Security team to enhance malware detection and characterization.",
		Tags:        []string{"AI/ML", "Security", "Big Data", "Android"},
	},
	{
		Company:     "SFEIR",
		Logo:        "sfeir.webp",
		Title:       "Data Engineer",
		BrandColor:  "#000000",
		Description: "Delivered data engineering, technical consulting, recruitment, and pre-sales.",
		Tags:        []string{"Data Engineering", "Consulting", "GCP"},
	},
	{
		Company:     "University of Luxembourg",
		Logo:        "uni-lu.webp",
		Title:       "Researcher & Teacher",
		BrandColor:  "#E3000F",
		Description: "Conducted AI security research and taught AI/ML, Big Data, and Android development.",
		Tags:        []string{"AI", "Security", "Research", "Teaching"},
	},
	{
		Company:     "Clearstream",
		Logo:        "clearstream.webp",
		Title:       "Security Engineer",
		BrandColor:  "#008C95",
		Description: "Selected and configured a scalable SIEM solution to detect and respond to security incidents.",
		Tags:        []string{"Big Data", "Security", "Banking", "SIEM"},
	},
}

var OPEN_SOURCE = []Project{
	{
		Title:       "mlops-python-package",
		Href:        "https://github.com/fmind/mlops-python-package",
		Description: "Kickstart your MLOps initiative with a flexible, robust, and productive Python package.",
	},
	{
		Title:       "cookiecutter-mlops-package",
		Href:        "https://github.com/fmind/cookiecutter-mlops-package",
		Description: "Start building and deploying Python packages and Docker images for MLOps.",
	},
	{
		Title:       "MLOps Coding Course",
		Href:        "https://mlops-coding-course.fmind.dev/",
		Repo:        "https://github.com/MLOps-Courses/mlops-coding-course",
		Description: "Learn to create, develop, and maintain a state-of-the-art MLOps code base.",
	},
}

var YOUTUBE_SERIES = []Playlist{
	{
		Title:       "Bleeding Agent",
		URL:         "https://www.youtube.com/playlist?list=PLPCnNL6Y2PbTckW80gDLnznFMDEz18HBS",
		Description: "Technical deep dives into the Black Box of AI Agents and emerging autonomous systems.",
		CTA:         "View Podcast",
	},
	{
		Title:       "AI Agents in a Nut$SHELL",
		URL:         "https://www.youtube.com/playlist?list=PLPCnNL6Y2PbT1aKOx2fMFBpicTRzeMS6f",
		Description: "Brief, high-signal deep dives into the core architecture and inner workings of AI Agents.",
		CTA:         "View Playlist",
	},
	{
		Title:       "MLOps Coding Course",
		URL:         "https://www.youtube.com/playlist?list=PLPCnNL6Y2PbQplCczUFhtQpCznEXqDZnh",
		Description: "Bridge the gap between robust software engineering and cutting-edge data science.",
		CTA:         "View Course",
	},
}

var THESIS = Thesis{
	Title:              "Creating better ground truth to further understand Android malware",
	URL:                "https://orbilu.uni.lu/handle/10993/39903",
	InstitutionDetails: "University of Luxembourg (SNT) & Google, 2019",
	Description:        "AI/ML models are only as good as the data they learn from — yet Android malware ground truths are notoriously unreliable. This thesis tackles the problem by benchmarking antivirus engines, harmonizing their conflicting labels, and mining large-scale datasets to characterize malicious behavior.",
	Links: []ThesisLink{
		{
			Label: "Servalx — Android Malware Processing",
			URL:   "https://github.com/fmind/servalx",
		},
		{
			Label: "APKWorkers — Distributed Android Analysis",
			URL:   "https://github.com/fmind/apkworkers",
		},
	},
}

var PAPERS = []ResearchPaper{
	{
		Title:     "Euphony: Harmonious Unification of Cacophonous Anti-Virus Vendor Labels",
		URL:       "https://orbilu.uni.lu/handle/10993/31441",
		Venue:     "MSR 2017 • Mining Software Repositories",
		Code:      "https://github.com/fmind/euphony",
		CodeLabel: "Euphony — Label Unification",
	},
	{
		Title:     "On the Lack of Consensus in Anti-Virus Decisions",
		URL:       "https://orbilu.uni.lu/handle/10993/27845",
		Venue:     "DIMVA 2016 • Detection of Intrusions and Malware",
		Code:      "https://github.com/fmind/stase",
		CodeLabel: "STASE — Statistical Metrics",
	},
}

func GetServices() []Service {
	return []Service{
		{
			Icon:        "🏢",
			Title:       "AI Architecture & Advisory",
			Description: "Engage me to assess, design, and scale AI initiatives — from agent platforms and MLOps to security and operating models.",
			Badge:       "🔴 Not available for new missions",
			BadgeType:   "error",
			CTAText:     "✉️ Get in Touch",
			CTAURL:      "mailto:" + METADATA.Email,
		},
		{
			Icon:        "🎓",
			Title:       "Mentoring",
			Description: "Book a paid 1-hour session to discuss your projects: upskilling, career mentoring, architecture review, brainstorming, and more.",
			Badge:       "💰 Paid session — 1 hour",
			BadgeType:   "info",
			CTAText:     "📅 Book a Session",
			CTAURL:      METADATA.CalendarURL,
		},
	}
}
