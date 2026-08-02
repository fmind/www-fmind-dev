package site

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// TestGuessLanguage pins the classifier against the shapes the archive actually
// contains. The imported articles carry no fenced languages, so this guess is
// the only thing standing between a code block and a wall of one color — and a
// wrong guess is worse than none, hence the "stays plain" cases.
func TestGuessLanguage(t *testing.T) {
	cases := []struct {
		name string
		code string
		want string
	}{
		{
			name: "python source",
			code: "import abc\nimport typing as T\n\ndef split(data):\n    return data\n",
			want: "python",
		},
		{
			name: "python without imports",
			code: "spaces = []\npage_token = None\nwhile True:\n    spaces.append(1)\n",
			want: "python",
		},
		{
			name: "agent config yaml quoting python",
			code: "# yaml-language-server: $schema=https://example/AgentConfig.json\nagent_class: LlmAgent\nmodel: gemini-2.5-flash\ninstruction: |\n  import this\n",
			want: "yaml",
		},
		{
			name: "http exchange",
			code: "POST /mcp HTTP/1.1\nMcp-Protocol-Version: 2026-07-28\n",
			want: "http",
		},
		{
			name: "json payload",
			code: "{\n  \"jsonrpc\": \"2.0\",\n  \"id\": 1\n}\n",
			want: "json",
		},
		{
			name: "shell commands",
			code: "pip install bromate\nuv run bromate --help\n",
			want: "bash",
		},
		{
			name: "shell session",
			code: "$ mise run test\nDONE 65 tests\n",
			want: "console",
		},
		{
			name: "dockerfile",
			code: "FROM python:3.12\nRUN pip install uv\n",
			want: "dockerfile",
		},
		{
			// Regression: both of this article's blocks fell through to plain text.
			name: "async def and decorated handler",
			code: "@app.post(\"/\")\nasync def index(request: Request) -> dict:\n    return {}\n",
			want: "python",
		},
		{
			name: "assignment from a call",
			code: "# main.py\nclient = genai.Client(\n    vertexai=True,\n)\n",
			want: "python",
		},
		{
			name: "http request without a version",
			code: "POST /api/agents/feature_request_agent/tasks\nContent-Type: application/json\n",
			want: "http",
		},
		{
			name: "command with a long flag",
			code: "adk web --a2a --port 8001\n",
			want: "bash",
		},
		{
			name: "argparse stays python",
			code: "parser = argparse.ArgumentParser()\nparser.add_argument(\"--verbose\")\n",
			want: "python",
		},
		{
			name: "directory tree stays plain",
			code: ".agents/docs/\n├── INDEX.md\n└── cloud-run/\n    └── DOC.md\n",
			want: "",
		},
		{
			name: "chat transcript stays plain",
			code: "Fred: So... what's on the agenda today?\nJamy: Something super smart!\n",
			want: "",
		},
		{
			name: "slash commands stay plain",
			code: "/plugin marketplace add fmind/agent-docs\n/plugin install agent-docs@agent-docs\n",
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := guessLanguage(tc.code); got != tc.want {
				t.Errorf("guessLanguage = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestHighlightRendersThemedTokens proves the pipeline end to end: a code block
// becomes class-tagged markup (never inline styles, which the CSP forbids) and
// the stylesheet that colors those classes comes from the same theme.
func TestHighlightRendersThemedTokens(t *testing.T) {
	if highlighterErr != nil {
		t.Fatalf("build highlighter: %v", highlighterErr)
	}
	article, err := parseArticle("content/articles/example.md", []byte(`+++
title = "Example"
description = "Example"
date = "2026-08-01"
slug = "example"
tags = ["Agent"]
+++

Body.

`+"```python"+`
def hello() -> str:
    return "world"
`+"```"+`
`), exampleAssets(t, ".webp"))
	if err != nil {
		t.Fatalf("parse article: %v", err)
	}
	if !strings.Contains(article.HTML, `<pre class="chroma">`) {
		t.Errorf("code block was not highlighted: %s", article.HTML)
	}
	if !strings.Contains(article.HTML, `<span class="k">def</span>`) {
		t.Errorf("keyword token is missing its class: %s", article.HTML)
	}
	if strings.Contains(article.HTML, "style=") {
		t.Error("highlighted markup uses inline styles, which the CSP blocks")
	}

	css, err := highlighter.CSS()
	if err != nil {
		t.Fatalf("generate stylesheet: %v", err)
	}
	if !strings.Contains(css, ".chroma .k") {
		t.Errorf("stylesheet does not color keyword tokens: %s", css)
	}
	if strings.Contains(css, "\n") {
		t.Error("generated stylesheet is not minified")
	}
}

// TestArchiveCodeBlocksAreHighlighted guards the wiring across the real archive,
// not a fixture: if the renderer is ever unregistered or the guesser regresses,
// every article silently falls back to one flat color and no other test notices.
// The floor is a ratio rather than "all blocks" on purpose — directory trees,
// console output, and chat transcripts have nothing to color and must stay plain.
func TestArchiveCodeBlocksAreHighlighted(t *testing.T) {
	const minimumHighlightedRatio = 0.75

	collection, err := loadArticles()
	if err != nil {
		t.Fatalf("load articles: %v", err)
	}
	blockPattern := regexp.MustCompile(`(?s)<pre class="chroma">.*?</pre>`)
	spanPattern := regexp.MustCompile(`<span class="([a-z0-9]+)"`)

	blocks, highlighted := 0, 0
	for _, article := range collection.all {
		for _, block := range blockPattern.FindAllString(article.HTML, -1) {
			blocks++
			for _, span := range spanPattern.FindAllStringSubmatch(block, -1) {
				// "line", "cl" and "w" are the formatter's structural wrappers; any
				// other class is a real token carrying a color.
				if class := span[1]; class != "line" && class != "cl" && class != "w" {
					highlighted++
					break
				}
			}
		}
	}
	if blocks == 0 {
		t.Fatal("the archive rendered no code blocks at all")
	}
	if ratio := float64(highlighted) / float64(blocks); ratio < minimumHighlightedRatio {
		t.Errorf("highlighted code blocks = %d/%d (%.0f%%), want at least %.0f%%",
			highlighted, blocks, ratio*100, minimumHighlightedRatio*100)
	}
}

// relativeLuminance implements the WCAG 2.1 definition so the contrast floor is
// checked against the spec rather than against eyeballed hex values.
func relativeLuminance(t *testing.T, value string) float64 {
	t.Helper()
	value = strings.TrimPrefix(value, "#")
	if len(value) == 3 {
		value = string([]byte{value[0], value[0], value[1], value[1], value[2], value[2]})
	}
	if len(value) != 6 {
		t.Fatalf("unexpected hex value %q", value)
	}
	weights := [3]float64{0.2126, 0.7152, 0.0722}
	luminance := 0.0
	for index := range 3 {
		parsed, err := strconv.ParseUint(value[index*2:index*2+2], 16, 8)
		if err != nil {
			t.Fatalf("parse hex value %q: %v", value, err)
		}
		channel := float64(parsed) / 255
		if channel <= 0.04045 {
			channel /= 12.92
		} else {
			channel = math.Pow((channel+0.055)/1.055, 2.4)
		}
		luminance += weights[index] * channel
	}
	return luminance
}

func contrastRatio(t *testing.T, foreground, background string) float64 {
	t.Helper()
	first, second := relativeLuminance(t, foreground), relativeLuminance(t, background)
	lighter, darker := math.Max(first, second), math.Min(first, second)
	return (lighter + 0.05) / (darker + 0.05)
}

var (
	// Chroma writes token rules as ".chroma .err" and wrapper rules as ".bg", so the
	// prefix is optional.
	cssRulePattern       = regexp.MustCompile(`/\* ([^*]+) \*/ (?:\.chroma )?\.([\w-]+) \{([^}]*)\}`)
	cssForegroundPattern = regexp.MustCompile(`(?:^|[;{ ])color:\s*(#[0-9a-fA-F]{3,6})`)
	cssBackgroundPattern = regexp.MustCompile(`background-color:\s*(#[0-9a-fA-F]{3,6})`)
)

// Code comments carry the explanation a reader most needs from a snippet, and
// Tokyo Night renders them at 1.91:1 against its own background. This asserts on
// the generated stylesheet rather than the style model, so it measures exactly
// what the browser receives — and swapping codeTheme later cannot quietly
// reintroduce unreadable text.
func TestCodeThemeTokensMeetContrastFloor(t *testing.T) {
	const minimumContrast = 4.5

	highlighter, err := newCodeHighlighter()
	if err != nil {
		t.Fatalf("new code highlighter: %v", err)
	}
	stylesheet, err := highlighter.CSS()
	if err != nil {
		t.Fatalf("generate stylesheet: %v", err)
	}

	background := cssBackgroundPattern.FindStringSubmatch(stylesheet)
	if background == nil {
		t.Fatal("generated stylesheet declares no code background")
	}

	checked := 0
	for _, rule := range cssRulePattern.FindAllStringSubmatch(stylesheet, -1) {
		token, class, body := strings.TrimSpace(rule[1]), rule[2], rule[3]
		foreground := cssForegroundPattern.FindStringSubmatch(body)
		if foreground == nil {
			continue
		}
		checked++
		if ratio := contrastRatio(t, foreground[1], background[1]); ratio < minimumContrast {
			t.Errorf("%s (.%s) %s on %s = %.2f:1, want >= %.1f:1",
				token, class, foreground[1], background[1], ratio, minimumContrast)
		}
	}
	if checked < 20 {
		t.Fatalf("only %d token rules inspected; the stylesheet did not generate", checked)
	}
}
