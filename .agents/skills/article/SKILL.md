---
name: article
description: Publish, revise, or retire an article on www.fmind.dev — import from the publications repo, generate image derivatives, and satisfy every startup and test invariant. Use for any change under content/articles/.
license: MIT
metadata:
  author: Médéric HURIER (Fmind)
---

# Publish an Article

Every published article lives in `content/articles/<slug>.md` and is validated at startup, not at request time. A malformed article does not render badly — it stops the binary from booting. This skill is the order of operations that keeps that gate green.

## Where an Article Comes From

Drafts are written in the private publications repository and enter this repository through `pub export`; that repository carries its own authoring instructions. Once it records the live site URL, the private draft is deleted and **this repository owns the only published body**. A substantial regenerated revision must therefore start from the current site Markdown, never from the retired private draft.

## Workflow

1. **Import the draft** with `pub export` from the publications repository. It writes the Markdown with its TOML frontmatter and applies the ~2.4MP pixel budget to body images.

1. **Check the tags** against the closed vocabulary in `templates/tags.go`. Anything outside it fails startup. Adding a tag is three coupled edits:
   - a `Tag` entry in `templates/tags.go` (declaration order is display order everywhere),
   - a matching `[data-tag='…']` color rule in `static/css/input.css`,
   - at least one article actually carrying it — an unused tag also fails startup.

1. **Generate the image derivatives** and commit them:
   ```bash
   mise run build:images        # FORCE=1 regenerates every source
   ```
   The cover's derivative is required to boot; any other missing rung fails the archive test. The ladder is `templates.DerivativeWidths` — widening it means a `FORCE=1` regeneration of the whole archive.

1. **Run the gates**:
   ```bash
   mise run format
   mise run check               # includes check:typos over content/articles
   mise run test
   ```
   Real terms that trip the spell check belong in `typos.toml` as reviewed exceptions — never weaken the check itself.

1. **Check the outbound links** before shipping — this one needs network and is not part of the offline `check`:
   ```bash
   mise run check:links
   ```

## Invariants Worth Knowing Before You Edit

- **Code blocks**: fence with a language. Unlabeled blocks fall back to the `languageMarkers` guesser in `highlight.go`; highlighting happens once at startup from `codeTheme`, which also generates the stylesheet.
- **Figures**: a standalone image renders as a `<figure>` that breaks out of the text column to `--figure-max-width` and links to full resolution. `figureSizes` in `articles.go` and the `.article-page` rules in `static/css/input.css` describe one layout and must change together.
- **Captions**: the paragraph after a standalone image folds into the figure as a `<figcaption>` when its text repeats the alt (`foldBodyCaptions`), and the folded image then drops its `alt`. Compare as **text**, never as markup — rendering adds links and typographic spaces. Position alone is not a caption signal: every article opens with a cover followed by ordinary prose.
- **Diagrams that read too small** are laid out wrong at the source, and are fixed there — never compensated for in the layout. Re-render from the diagram source and re-import; the authoring rules for that live with the diagram sources in the publications repository. What this repository asserts is only the acceptance bar: labels ≥ ~12px apparent size and rendered height ≤ ~1300px at 1280 wide.
- **Drafts never leak**: the validated article collection is the single source for HTML, Atom, sitemap, `llms.txt`, JSON, and MCP surfaces, and production discovery excludes drafts.

## Gotchas

1. **Startup is the test**: most article mistakes surface as a boot failure, so `mise run test` catches them — do not wait for a deploy to find out.
2. **Derivatives are committed artifacts**: they are generated deterministically by the pinned pure-Go WebP encoder under the `nodynamic` build tag, so regenerating on any machine produces identical bytes. Commit them.
3. **Never upscale a rung**: candidates are discovered from disk, so an image offers only the rungs its source was wide enough to earn.
