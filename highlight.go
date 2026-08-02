package site

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// codeTheme is the Chroma style every code block is rendered with. Tokyo Night's
// palette sits next to the brand navy, so a code block reads as part of the page
// instead of a foreign panel. Any name from chroma/styles works here ("monokai",
// "github-dark", …); the markup and the generated CSS both follow it.
const codeTheme = "tokyonight-night"

// Tokyo Night ships two token groups that fail WCAG AA against its own background
// (#1a1b26): comments and docstrings at #414868 reach only 1.91:1, and the error
// red #db4b4b reaches 4.17:1. Comments carry the explanation a reader most needs
// from a snippet, so they are the worst possible place to lose legibility. These
// overrides keep the palette and lift both groups past 4.5:1 — #7c86b4 measures
// 4.82:1 and stays recessed against the code around it, and #f7768e (6.46:1) is
// Tokyo Night's own red rather than a new hue. Verify with a contrast checker
// before changing either value or codeTheme itself.
var codeThemeContrastFixes = chroma.StyleEntries{
	chroma.Comment:              "#7c86b4",
	chroma.CommentHashbang:      "#7c86b4",
	chroma.CommentMultiline:     "#7c86b4",
	chroma.CommentPreproc:       "#7c86b4",
	chroma.CommentPreprocFile:   "#7c86b4",
	chroma.CommentSingle:        "#7c86b4",
	chroma.CommentSpecial:       "#7c86b4",
	chroma.LiteralStringDoc:     "#7c86b4",
	chroma.LiteralStringHeredoc: "#7c86b4",
	chroma.Error:                "#f7768e",
	chroma.GenericDeleted:       "#f7768e",
	chroma.GenericError:         "#f7768e",
	chroma.GenericTraceback:     "#f7768e",
}

// codeHighlighter renders code blocks with Chroma at startup, so highlighted
// markup is part of the pre-rendered article and the page still ships no
// JavaScript. Classes rather than inline styles: the CSP allows no style
// attributes, and one stylesheet beats a color repeated on every token.
type codeHighlighter struct {
	formatter *chromahtml.Formatter
	style     *chroma.Style
}

func newCodeHighlighter() (*codeHighlighter, error) {
	style := styles.Get(codeTheme)
	// styles.Get falls back to a default instead of failing, which would silently
	// ship the wrong palette after a typo in codeTheme.
	if style == nil || !strings.EqualFold(style.Name, codeTheme) {
		return nil, fmt.Errorf("unknown code theme %q", codeTheme)
	}
	style, err := style.Builder().AddAll(codeThemeContrastFixes).Build()
	if err != nil {
		return nil, fmt.Errorf("apply %q contrast fixes: %w", codeTheme, err)
	}
	return &codeHighlighter{formatter: chromahtml.New(chromahtml.WithClasses(true)), style: style}, nil
}

// RegisterFuncs takes over both code block kinds. The imported archive predates
// this site and writes every snippet as an indented block, so handling fenced
// blocks alone would highlight nothing that is published today.
func (highlighter *codeHighlighter) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, highlighter.render)
	reg.Register(ast.KindCodeBlock, highlighter.render)
}

func (highlighter *codeHighlighter) render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	var code bytes.Buffer
	lines := node.Lines()
	for i := range lines.Len() {
		line := lines.At(i)
		code.Write(line.Value(source))
	}

	language := ""
	if fenced, ok := node.(*ast.FencedCodeBlock); ok && fenced.Info != nil {
		language = string(fenced.Language(source))
	}
	iterator, err := highlighter.lexer(language, code.String()).Tokenise(nil, code.String())
	if err != nil {
		return ast.WalkStop, fmt.Errorf("tokenise %q code block: %w", language, err)
	}
	if err := highlighter.formatter.Format(w, highlighter.style, iterator); err != nil {
		return ast.WalkStop, fmt.Errorf("format %q code block: %w", language, err)
	}
	// The formatter emits the whole <pre>, so goldmark must not also render the
	// block's raw text underneath it.
	return ast.WalkSkipChildren, nil
}

// lexer resolves the declared language, then the guess, then plain text.
func (highlighter *codeHighlighter) lexer(language, code string) chroma.Lexer {
	if language == "" {
		language = guessLanguage(code)
	}
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	// Coalesce merges runs of same-type tokens, which keeps the emitted markup
	// (and every rendered page) meaningfully smaller.
	return chroma.Coalesce(lexer)
}

// languageMarkers identify an unlabeled code block by a single high-signal
// pattern each, in decreasing order of confidence. Chroma's own analyser is not
// used: it scores weak matches, so it painted HTTP requests as arithmetic and
// left plain Python untouched. Order matters — YAML's "key: value" shape is the
// loosest signal here and must not win over a language that declared itself.
var languageMarkers = []struct {
	marker   *regexp.Regexp
	language string
}{
	// A request line, with the version or with a header line under it.
	{language: "http", marker: regexp.MustCompile(`(?m)^(GET|POST|PUT|PATCH|DELETE|HEAD) /\S*( HTTP/\d|\n[A-Za-z][\w-]*: )`)},
	{language: "html", marker: regexp.MustCompile(`(?m)^\s*<(!DOCTYPE|html|head|body|script|template|div)\b`)},
	{language: "dockerfile", marker: regexp.MustCompile(`(?m)^FROM \S+`)},
	{language: "json", marker: regexp.MustCompile(`\A\s*[{\[][\s\S]*"\s*:`)},
	// A document that opens on top-level keys is YAML even when a value further
	// down quotes Python: config files embed code, source files do not open on keys.
	{language: "yaml", marker: regexp.MustCompile(`(?m)\A(?:#[^\n]*\n)*[a-z][\w.-]*:(\s|$)[\s\S]*^[a-z][\w.-]*:(\s|$)`)},
	// Python, the archive's dominant language, gets several markers because its
	// snippets rarely open on an import: most are excerpts that start mid-module.
	// Definitions and imports, including `async def`.
	{language: "python", marker: regexp.MustCompile(`(?m)^\s*((async )?def \w+\s*\(|class \w+[\s(:]|from \S+ import |import \w+)`)},
	// Decorators, with or without arguments (`@dataclass`, `@app.post("/")`).
	{language: "python", marker: regexp.MustCompile(`(?m)^\s*@\w[\w.]*(\(|\s*$)`)},
	// An annotated return type: no other language in this archive writes `-> T:`.
	{language: "python", marker: regexp.MustCompile(`(?m)\)\s*->\s*[\w\[\]., |]+:[ \t]*$`)},
	// Block openers only Python spells this way. "if"/"with"/"else" are left out on
	// purpose: they are also GitHub Actions YAML keys.
	{language: "python", marker: regexp.MustCompile(`(?m)^[ \t]*(for|while|elif|try|except|finally)\b[^\n]*:[ \t]*$`)},
	// Assignment from a call — `client = genai.Client(`. Weak on its own, which is
	// why it runs after every self-declaring format above has had its turn.
	{language: "python", marker: regexp.MustCompile(`(?m)^\s*[A-Za-z_]\w* = [A-Za-z_][\w.]*\(`)},
	{language: "go", marker: regexp.MustCompile(`(?m)^(package \w+|func \w+\()`)},
	{language: "sql", marker: regexp.MustCompile(`(?im)^SELECT\s[\s\S]*\sFROM\s`)},
	{language: "toml", marker: regexp.MustCompile(`(?m)^\[[\w.-]+\]$[\s\S]*^[\w.-]+ = `)},
	// A run of bare `key = value` lines with no section header: a TOML excerpt.
	// Three in a row, because two could be consecutive Python assignments.
	{language: "toml", marker: regexp.MustCompile(`(?m)^[\w.-]+ = [^\n(]+\n[\w.-]+ = [^\n(]+\n[\w.-]+ = `)},
	// A justfile: an attribute, or a recipe whose indented body interpolates
	// `{{ }}`. Chroma maps the "justfile" name onto its Makefile lexer. This runs
	// before the shell markers because recipe bodies are full of shell commands.
	{language: "justfile", marker: regexp.MustCompile(`(?m)^\[[a-z-]+[(\]]`)},
	{language: "justfile", marker: regexp.MustCompile(`(?m)^[\w-]+( \w+="[^"]*")*:[ \t]*$\n[ \t]+\S[\s\S]*\{\{`)},
	{language: "console", marker: regexp.MustCompile(`(?m)^\$ \S`)},
	{language: "bash", marker: regexp.MustCompile(`(?m)^(pip|uv|docker|kubectl|gcloud|git|curl|mise|npm|npx|poetry|python|make|export|cd|sudo|terraform) \S`)},
	// Any command invoked with a long flag ("adk web --a2a", "pandoc --to=plain").
	// Python is matched above, so an argparse line cannot reach this rule.
	{language: "bash", marker: regexp.MustCompile(`(?m)^[a-z][\w.-]+[^\n]* --[\w-]`)},
	// Config excerpts rarely start at the top of their file: they open on a list
	// item ("- run: …") or on a key whose value is the indented block below it.
	{language: "yaml", marker: regexp.MustCompile(`(?m)^[ \t]*-[ \t]+[a-z][\w.-]*:(\s|$)`)},
	{language: "yaml", marker: regexp.MustCompile(`(?m)^[a-z][\w.-]*:[ \t]*(#[^\n]*)?$\n[ \t]+\S`)},
	// A pinned dependency file (requirements.txt, constraints.txt). It opens on a
	// "# comment" like a shell script and carries "# via" lines like a heading, so
	// it must be claimed before the shell and Markdown markers see it.
	{language: "text", marker: regexp.MustCompile(`(?m)^[\w.-]+==\d`)},
	// A document quoted as a code block: a heading plus at least one more block
	// structure. One heading alone is a shell comment far more often than prose.
	{language: "markdown", marker: regexp.MustCompile(`(?m)\A#{1,6} \S[\s\S]*^(#{1,6} |[ \t]*[-*] |\|)`)},
	// A table needs no heading above it to be unmistakably Markdown.
	{language: "markdown", marker: regexp.MustCompile(`\A\|[^\n]*\|[ \t]*\n[ \t]*\|[-: |]+\|`)},
	// Lowercase keys only: a capitalised "Name:" line is prose or a transcript far
	// more often than it is configuration.
	{language: "yaml", marker: regexp.MustCompile(`(?m)^[a-z][\w.-]*:(\s|$)[\s\S]*^[a-z][\w.-]*:(\s|$)`)},
}

// guessLanguage names the lexer for a code block that declared none. It returns
// an empty string unless a marker is unambiguous: a block rendered in the wrong
// palette reads worse than one rendered in none.
func guessLanguage(code string) string {
	for _, candidate := range languageMarkers {
		if candidate.marker.MatchString(code) {
			return candidate.language
		}
	}
	return ""
}

// CSS returns the stylesheet for the theme's token classes. It is appended to
// the inlined stylesheet at startup rather than committed as a generated file,
// so the palette can never drift from the markup it colors.
func (highlighter *codeHighlighter) CSS() (string, error) {
	var css bytes.Buffer
	if err := highlighter.formatter.WriteCSS(&css, highlighter.style); err != nil {
		return "", fmt.Errorf("write %q stylesheet: %w", codeTheme, err)
	}
	// Chroma writes one rule per line; the inlined stylesheet is minified, so the
	// newlines would be the only unminified bytes on every page.
	return strings.ReplaceAll(css.String(), "\n", ""), nil
}
