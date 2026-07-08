package site

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/fmind/www-fmind-dev/templates"
)

// mcpProfileURI is the stable URI of the full-portfolio MCP resource.
const mcpProfileURI = "portfolio://profile.json"

// noArgs is the input type for tools that take no parameters. The MCP SDK infers
// an empty-object JSON schema from it, which the spec requires for every tool.
type noArgs struct{}

// Portfolio is the complete, machine-readable snapshot of the site's content. It
// backs both the MCP `portfolio://profile.json` resource and the get_profile tool.
type Portfolio struct {
	Metadata        templates.Metadata             `json:"metadata"`
	Biography       []string                       `json:"biography"`
	Expertise       []templates.ExpertiseCard      `json:"expertise"`
	Experience      []templates.WorkExperience     `json:"experience"`
	Certifications  []templates.CertificationBadge `json:"certifications"`
	Specializations []templates.CertificationEntry `json:"specializations"`
	Thesis          templates.Thesis               `json:"thesis"`
	Papers          []templates.ResearchPaper      `json:"papers"`
	Posts           []templates.CuratedPost        `json:"posts"`
	OpenSource      []templates.Project            `json:"open_source"`
	YouTubeSeries   []templates.Playlist           `json:"youtube_series"`
	Services        []templates.Service            `json:"services"`
}

// snapshot assembles the current portfolio from the templates data package.
func snapshot() Portfolio {
	return Portfolio{
		Metadata:        templates.METADATA,
		Biography:       templates.BIOGRAPHY,
		Expertise:       templates.EXPERTISE,
		Experience:      templates.EXPERIENCES,
		Certifications:  templates.BADGES,
		Specializations: templates.SPECIALIZATIONS,
		Thesis:          templates.THESIS,
		Papers:          templates.PAPERS,
		Posts:           templates.POSTS,
		OpenSource:      templates.OPEN_SOURCE,
		YouTubeSeries:   templates.YOUTUBE_SERIES,
		Services:        templates.GetServices(),
	}
}

// profileResult is the get_profile tool's structured output: the "who am I" core.
type profileResult struct {
	Metadata  templates.Metadata        `json:"metadata"`
	Biography []string                  `json:"biography"`
	Expertise []templates.ExpertiseCard `json:"expertise"`
}

type experienceResult struct {
	Experience []templates.WorkExperience `json:"experience"`
}

type certificationsResult struct {
	Certifications  []templates.CertificationBadge `json:"certifications"`
	Specializations []templates.CertificationEntry `json:"specializations"`
}

type publicationsResult struct {
	Thesis templates.Thesis          `json:"thesis"`
	Papers []templates.ResearchPaper `json:"papers"`
	Posts  []templates.CuratedPost   `json:"posts"`
}

type projectsResult struct {
	OpenSource    []templates.Project  `json:"open_source"`
	YouTubeSeries []templates.Playlist `json:"youtube_series"`
}

type servicesResult struct {
	Services []templates.Service `json:"services"`
}

// readOnly marks every tool as non-mutating so MCP clients (Claude, ChatGPT) skip
// write-confirmation friction — the whole surface is read-only by construction.
var readOnly = &mcp.ToolAnnotations{ReadOnlyHint: true}

// newMCPServer builds the read-only MCP server exposing the portfolio as agent
// tools plus a single JSON resource. All handlers are pure reads of static data.
func newMCPServer() *mcp.Server {
	server := mcp.NewServer(
		&mcp.Implementation{Name: "www-fmind-dev", Version: "1.0.0", Title: "Médéric Hurier (Fmind) — Portfolio"},
		&mcp.ServerOptions{Instructions: "Query the portfolio of Médéric Hurier (Fmind), freelance AI/ML Architect & Engineer: profile, work experience, certifications, publications, projects, and services."},
	)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_profile",
		Description: "Return the core profile: identity, headline, job title, contact, socials, biography, and areas of expertise.",
		Annotations: readOnly,
	}, func(context.Context, *mcp.CallToolRequest, noArgs) (*mcp.CallToolResult, profileResult, error) {
		return nil, profileResult{Metadata: templates.METADATA, Biography: templates.BIOGRAPHY, Expertise: templates.EXPERTISE}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_experience",
		Description: "List professional work experience: companies, roles, descriptions, and skill tags.",
		Annotations: readOnly,
	}, func(context.Context, *mcp.CallToolRequest, noArgs) (*mcp.CallToolResult, experienceResult, error) {
		return nil, experienceResult{Experience: templates.EXPERIENCES}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_certifications",
		Description: "List professional certifications, badges, and course specializations with their issuers.",
		Annotations: readOnly,
	}, func(context.Context, *mcp.CallToolRequest, noArgs) (*mcp.CallToolResult, certificationsResult, error) {
		return nil, certificationsResult{Certifications: templates.BADGES, Specializations: templates.SPECIALIZATIONS}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_publications",
		Description: "List academic and written publications: the PhD thesis, peer-reviewed papers, and curated blog posts.",
		Annotations: readOnly,
	}, func(context.Context, *mcp.CallToolRequest, noArgs) (*mcp.CallToolResult, publicationsResult, error) {
		return nil, publicationsResult{Thesis: templates.THESIS, Papers: templates.PAPERS, Posts: templates.POSTS}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_projects",
		Description: "List open-source projects and educational YouTube series.",
		Annotations: readOnly,
	}, func(context.Context, *mcp.CallToolRequest, noArgs) (*mcp.CallToolResult, projectsResult, error) {
		return nil, projectsResult{OpenSource: templates.OPEN_SOURCE, YouTubeSeries: templates.YOUTUBE_SERIES}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_services",
		Description: "List the professional services on offer (freelancing, mentoring) with availability and how to book.",
		Annotations: readOnly,
	}, func(context.Context, *mcp.CallToolRequest, noArgs) (*mcp.CallToolResult, servicesResult, error) {
		return nil, servicesResult{Services: templates.GetServices()}, nil
	})

	server.AddResource(&mcp.Resource{
		URI:         mcpProfileURI,
		Name:        "profile",
		Title:       "Full portfolio",
		Description: "The complete portfolio (profile, experience, certifications, publications, projects, services) as a single JSON document.",
		MIMEType:    "application/json",
	}, func(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != mcpProfileURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}
		data, err := json.Marshal(snapshot())
		if err != nil {
			return nil, fmt.Errorf("marshaling portfolio: %w", err)
		}
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{{URI: req.Params.URI, MIMEType: "application/json", Text: string(data)}},
		}, nil
	})

	return server
}

// newMCPHandler exposes the MCP server over the stateless Streamable HTTP
// transport as a plain http.Handler, mounted by the router at /mcp. JSONResponse
// keeps responses as application/json (no SSE), which is simplest behind the
// Gateway proxy for a read-only, horizontally-scaled service.
func newMCPHandler() http.Handler {
	server := newMCPServer()
	return mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return server },
		&mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true},
	)
}
