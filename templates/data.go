package templates

import (
	"fmt"
	"time"
)

// Data Models
type SocialLink struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Icon string `json:"icon"`
}

type Metadata struct {
	Name              string       `json:"name"`
	AlternateName     string       `json:"alternate_name"`
	Title             string       `json:"title"`
	JobTitle          string       `json:"job_title"`
	HeadlinePrimary   string       `json:"headline_primary"`
	HeadlineSecondary string       `json:"headline_secondary"`
	Description       string       `json:"description"`
	Keywords          []string     `json:"keywords"`
	Email             string       `json:"email"`
	CalendarURL       string       `json:"calendar_url"`
	SiteURL           string       `json:"site_url"`
	TwitterHandle     string       `json:"twitter_handle"`
	Socials           []SocialLink `json:"socials"`
}

type CertificationBadge struct {
	URL    string `json:"url"`
	Logo   string `json:"logo"`
	Title  string `json:"title"`
	Issuer string `json:"issuer"`
	CertID string `json:"cert_id"`
	Active bool   `json:"active"`
}

type CertificationEntry struct {
	URL           string `json:"url"`
	Logo          string `json:"logo"`
	Title         string `json:"title"`
	IssuerDetails string `json:"issuer_details"`
}

type WorkExperience struct {
	Company     string   `json:"company"`
	Logo        string   `json:"logo"`
	Title       string   `json:"title"`
	BrandColor  string   `json:"brand_color"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type Project struct {
	Title       string `json:"title"`
	Href        string `json:"href"`
	Repo        string `json:"repo,omitempty"`
	Description string `json:"description"`
}

type Playlist struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	CTA         string `json:"cta"`
}

type ThesisLink struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type Thesis struct {
	Title              string       `json:"title"`
	URL                string       `json:"url"`
	InstitutionDetails string       `json:"institution_details"`
	Description        string       `json:"description"`
	Links              []ThesisLink `json:"links"`
}

type ResearchPaper struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Venue     string `json:"venue"`
	Code      string `json:"code"`
	CodeLabel string `json:"code_label"`
}

type CuratedPost struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Service struct {
	Icon        string `json:"icon"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Badge       string `json:"badge"`
	BadgeType   string `json:"badge_type"`
	CTAText     string `json:"cta_text"`
	CTAURL      string `json:"cta_url"`
}

type ExpertiseCard struct {
	Title       string `json:"title"`
	Emoji       string `json:"emoji"`
	Description string `json:"description"`
}

// Portfolio Data
var METADATA = Metadata{
	Name:              "Médéric Hurier",
	AlternateName:     "Fmind",
	Title:             "Médéric Hurier (Fmind) | AI/ML Architect & Engineer",
	JobTitle:          "AI/ML Architect & Engineer",
	HeadlinePrimary:   "Freelancer • AI/ML Architect & Engineer • AI Agents & MLOps",
	HeadlineSecondary: "GCP Professional Cloud Architect • PhD in AI & Computer Security",
	Description:       "Freelance AI/ML Architect & Engineer specializing in AI Agents, Generative AI, MLOps, and Google Cloud Platform. PhD in AI & Computer Security.",
	Keywords: []string{
		"AI",
		"Machine Learning",
		"MLOps",
		"Artificial Intelligence",
		"AI Agents",
		"Agentic AI",
		"Generative AI",
		"Google Cloud",
		"GCP",
		"Python",
		"Freelance",
		"Luxembourg",
	},
	Email:         "contact@fmind.dev",
	CalendarURL:   "https://calendar.google.com/calendar/u/0/appointments/schedules/AcZssZ2ye3X9589PA2xmbV73Iz5J_NbFig6nN651vn6UuYAC-Cs5vBxnQ2L5db9UnAeXmUBQSW1MOobd",
	SiteURL:       "https://www.fmind.dev",
	TwitterHandle: "@fmind_dev",
	Socials: []SocialLink{
		{Name: "LinkedIn", URL: "https://www.linkedin.com/in/fmind-dev/", Icon: "linkedin"},
		{Name: "X (Twitter)", URL: "https://x.com/fmind_dev", Icon: "x"},
		{Name: "Medium", URL: "https://fmind.medium.com/", Icon: "medium"},
		{Name: "GitHub", URL: "https://github.com/fmind", Icon: "github"},
		{Name: "YouTube", URL: "https://www.youtube.com/@fmind-dev", Icon: "youtube"},
	},
}

var BIOGRAPHY = []string{
	"As a freelance <strong>AI/ML Architect &amp; Engineer</strong> with a Ph.D. in AI and Computer Security, I specialize in designing autonomous agentic workflows and production-ready MLOps environments. My focus is on driving real-world business outcomes by ensuring reliability and observability in complex systems.",
	"Certified as a <strong>Google Cloud Professional Cloud Architect</strong>, I combine deep academic research with a pragmatic, cartesian approach to engineering. I have helped institutions like Google, Decathlon, BNP Paribas, ArcelorMittal, and the European Commission grow their AI/ML projects to enterprise-grade solutions in production.",
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
		Description: "Leveraging Ph.D. background to build secure and trustworthy AI systems.",
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
		Title:       "AI Agents in a Nut$HELL",
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
			Label: "SetValX — Android Malware Processing",
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

var POSTS = []CuratedPost{
	{
		Title: "Architecting the AI Agent Platform: A Definitive Guide",
		URL:   "https://fmind.medium.com/architecting-the-ai-agent-platform-a-definitive-guide-405750a3de44",
	},
	{
		Title: "Powering Up Your Agent in Production with ADK, OAuth, and Gemini Enterprise",
		URL:   "https://fmind.medium.com/powering-up-your-agent-in-production-with-adk-oauth-and-gemini-enterprise-a52b0716fcba",
	},
	{
		Title: "DA2A: The Future of Data Platforms Is Agentic, Distributed and Collaborative",
		URL:   "https://fmind.medium.com/da2a-the-future-of-data-platforms-is-agentic-distributed-and-collaborative-741d5aa96fc4",
	},
	{
		Title: "CAG vs RAG: Choosing the Right Strategy for Your AI Application",
		URL:   "https://fmind.medium.com/cag-vs-rag-choosing-the-right-strategy-for-your-ai-application-68dcae85d028",
	},
	{
		Title: "Stop Building Rigid AI/ML Pipelines: Embrace Reusable Components",
		URL:   "https://fmind.medium.com/stop-building-rigid-ai-ml-pipelines-embrace-reusable-components-for-flexible-mlops-6e165d837110",
	},
	{
		Title: "Poetry Was Good, uv Is Better: An MLOps Migration Story",
		URL:   "https://fmind.medium.com/poetry-was-good-uv-is-better-an-mlops-migration-story-f52bf0c6c703",
	},
	{
		Title: "Make Your MLOps Code Base SOLID with Pydantic and Python's ABC",
		URL:   "https://fmind.medium.com/make-your-mlops-code-base-solid-with-pydantic-and-pythons-abc-aeedfe9c3e65",
	},
	{
		Title: "How to Configure VS Code for AI/ML and MLOps Development in Python 🐍",
		URL:   "https://fmind.medium.com/how-to-configure-vs-code-for-ai-ml-and-mlops-development-in-python-%EF%B8%8F%EF%B8%8F-8582d8c6ea54",
	},
}

func GetServices() []Service {
	year := time.Now().Year()
	return []Service{
		{
			Icon:        "🏢",
			Title:       "Freelancing",
			Description: "Hire me as a freelance AI/ML Architect & Engineer to design, build, and scale your AI initiatives — from AI Agents to MLOps.",
			Badge:       fmt.Sprintf("⏳ Not available until end of %d", year),
			BadgeType:   "warning",
			CTAText:     "✉️ Contact Me",
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
