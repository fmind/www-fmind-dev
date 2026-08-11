+++
title = "The MLOps Adventure Continues: An AGENTS.md Ready Stack for AI/ML"
description = "AI coding agents exposed the cost and drift of a fragmented Python toolchain. This guide moves formatting, checking, testing, and releases into one AGENTS.md-ready stack."
date = "2026-08-11"
tags = ["Coding", "MLOps", "Guide"]
slug = "mlops-adventure-continue"
draft = false
+++

![The MLOps Adventure Continues: An AGENTS.md Ready Stack for AI/ML](/static/img/articles/mlops-adventure-continue/cover.webp)

An MLOps stack is never finished. The stack behind my [MLOps Coding Course](https://mlops-coding-course.fmind.dev) tells thousands of engineers "this is the state-of-the-art way to set up an ML project", which makes every pinned tool a promise with a shelf life. The version I built two years ago ran on Python 3.13, MLflow 2, mypy, `just`, and `pre-commit`. By mid-2026 the ground had moved: Python 3.14 is stable, MLflow is on 3.x, and the [Astral](https://astral.sh) ecosystem behind [uv](https://docs.astral.sh/uv/) and [Ruff](https://docs.astral.sh/ruff/) grew a type checker, [ty](https://github.com/astral-sh/ty).

But the real reason to move was not tool churn. AI coding agents are the first contributors that read your entire repository and run your commands on day one - and they amplify whatever they find. Hand an agent a fuzzy toolchain and it will confidently multiply the confusion. Hand it one deterministic gate and machine-readable rules, and its first attempt lands close. Good tools are not a convenience for agents; they are the guardrails that make autonomy safe.

So twice this summer, the four repositories behind the course shipped coordinated major versions: the July round built the guardrails, the August round proved they work. The headline number: both majors are net deletions - 532 more lines removed than added in `v5.0.0`, another 393 in `v6.0.0` - while the gate grew to eight static checks and test coverage rose from 80 to 100 percent. Less code, more proof. That is what technical debt leaving the building looks like.

## What Two Majors Bought

Four gains, each with evidence behind it:

- **Robustness.** Coverage went from 80 to 100 percent, and CI now proves what it claims: it runs the same gate as the local hooks, then asserts the working tree is untouched. Before, nothing but my laptop had ever confirmed the package builds a wheel.
- **Speed.** One command answers "is this correct?" identically everywhere - locally, in CI, or under an agent - so iterations go into the problem instead of into environment drift. The Rust-based tools keep each check fast even as the gate grew to eight of them.
- **Security.** A routine check caught 25 known vulnerabilities with zero commits of my own, and three new linters now cover the workflows and Dockerfile that nothing checked before.
- **Best practices, enforced.** The course finally cuts its own semver releases, and every rule an agent needs lives in one `AGENTS.md`, verified against the code it describes.

The rest of this article is the evidence.

## Four Repositories, One Promise

The contradiction that started it all: the course has a complete chapter on cutting releases with semantic versioning and a changelog, and until July its own repository had zero git tags. It taught a discipline it never practiced.

The promise spans four repositories that only work when they agree:

- **[mlops-python-package](https://github.com/fmind/mlops-python-package/releases/tag/v6.0.0)** (`v6.0.0`) - the reference implementation: a production-shaped ML pipeline driven by typed configs, MLflow, and [Pydantic](https://docs.pydantic.dev) + [Pandera](https://pandera.readthedocs.io) validation.
- **[cookiecutter-mlops-package](https://github.com/fmind/cookiecutter-mlops-package/releases/tag/v6.0.0)** (`v6.0.0`) - the template that scaffolds a new project shaped like the package above.
- **[mlops-coding-course](https://github.com/MLOps-Courses/mlops-coding-course/releases/tag/v7.0.0)** (`v7.0.0`) - the free course that explains every decision, chapter by chapter.
- **[mlops-coding-skills](https://github.com/MLOps-Courses/mlops-coding-skills/releases/tag/v2.0.0)** (`v2.0.0`) - the same methodology packaged as seven Agent Skills for AI coding assistants.

Drift in one silently contradicts the other three: the template scaffolds commands the course no longer teaches, or a skill tells an agent to run a type checker the package has dropped. That coupling is why the releases ship in lockstep - and why a shared toolchain pays twice. After the August release, the same one-line CI fix landed in three of the four repositories on the same day.

## One Gate, Named Once

The July round rebuilt the stack around a single idea: one shared vocabulary that humans, git hooks, CI, and agents all call. `just` gave way to [mise](https://mise.jdx.dev), `pre-commit` to [lefthook](https://github.com/evilmartians/lefthook), mypy to `ty`, bandit to Ruff's `S` rules, `commitizen` to [git-cliff](https://git-cliff.org), hatchling to `uv_build`, and [dprint](https://dprint.dev) took over config and markup formatting. The through-line is not "Rust is faster" - it is that one coherent ecosystem removes the seams where tools disagree, and those seams are exactly what confuse a new contributor, human or not.

The August round finished the thought by naming the gate once:

```toml
[tasks.all]
alias = "a"
description = "Format, check, test, and build the project (the canonical gate)"
run = ["mise run format", "mise run check", "mise run test", "mise run build"]
```

CI now runs that single task and asserts the working tree is clean afterwards. It sounds cosmetic; it is a correctness property. When CI enumerated the steps itself, it was a second list, and second lists drift: mine had already lost `mise run build`, and the template's bake test ran every generated task except `mise run test` - so the template shipped a pytest suite its own CI had never executed.

A list of steps can omit one. A named gate cannot.

For an agent, this is the guardrail that matters most. An agent runs your gate far more often than you do, and a check that gives different verdicts locally and in CI burns its iterations on the disagreement. One command, one verdict, reproducible anywhere.

## Entropy With Zero Commits

A month after the July release, with `main` still byte-identical to the `v5.0.0` tag, `mise run check` on a clean checkout came back red: 25 known vulnerabilities across four transitive dependencies of MLflow. Nothing had changed except the world.

That is the claim worth internalizing: a dependency audit is not a checkbox you pass at release time. It is a clock running against your lockfile - and the stack rang the alarm before any reader hit the rot.

One fix separates two things that look alike. `cryptography` 50.0.0 patches PYSEC-2026-3552, a padding-oracle vulnerability - but MLflow declares `cryptography<50`, so the audit demanded a version the resolver refused to install. I overrode the stale bound with `override-dependencies = ["cryptography>=50"]` and proved it in both directions: the full suite and all six MLflow jobs pass with the override, and removing it re-locks to 49.0.0 and fails the audit again. Overriding a constraint and proving the result is engineering; silencing the scanner is a decision to ship the vulnerability quietly.

## What the New Checks Caught

Workflows and the Dockerfile were the only code in the repositories nothing linted, so the August round wired in [actionlint](https://github.com/rhysd/actionlint), [zizmor](https://docs.zizmor.sh), and [hadolint](https://github.com/hadolint/hadolint). The first runs paid for the setup:

- **A cache-poisoning pattern** in a release-triggered workflow that pushes container images - inherited by every project ever generated from the template.
- **A permissions bug** in the course's Pages deploy: a job-level `permissions:` block replaces the workflow-level one instead of merging with it, so the deploy job was silently dropping a grant.
- **A one-word deadlock no linter could see.** The template generated a CI job named `check` and a branch ruleset requiring `checks` - so every generated project that installed its ruleset had pull requests blocked forever, waiting on a check that could never report. The fix renames the job and writes the coupling into both `AGENTS.md` files, so the next contributor - human or agent - is told the names must match.
- **A security scan that was green because it was not scanning.** `check:scan` ran `trivy config .`, which enables only one of the four scanners declared in `trivy.yaml`. It had been green for a month, and green was the problem.

The same round deleted debt outright: a 790 KB SQLite database committed by accident, a changelog filter that had never matched Dependabot's commit prefix, and a formatter exclusion that left the skills repository's seven `SKILL.md` files - the product itself - as the only ones never formatted.

## Off the File Store: MLflow on SQLite

The one breaking change was overdue. MLflow 3 put the local file store in maintenance mode, and my July release handled it the lazy way: an opt-in flag plus a note saying "use a database in production". Reading the installed source showed MLflow had already moved on:

```python
DEFAULT_TRACKING_URI = "sqlite:///mlflow.db"
```

My package was not being dragged along by a deprecation; it was explicitly opting out of the new default. The fix is mostly deletion: `MlflowService` now defaults to `sqlite:///mlflow.db` for tracking and registry, and the opt-in flag is gone.

The honest cost: the test suite went from 34 to 339 seconds, because every test migrated a fresh SQLite store. A session-scoped fixture that migrates one template database and copies the file per test brings it back to roughly double the file store - a factor of two I will not pretend away. What it buys is a local setup with the same shape as production, a model registry on the store it was designed for, and a promotion to a real server that is one environment variable instead of a rewrite. I verified the whole path: all six pipeline jobs run on SQLite, register the model, and promote it to the `Champion` alias.

## Prose Does Not Run

Executable guardrails are half the story. The other half is `AGENTS.md` - an [open format stewarded by the Agentic AI Foundation](https://agents.md/) under the Linux Foundation - which records what a gate cannot: the exact commands, the definition of done, the conventions. One file, readable by any agent, replaces the fragmented pile of per-tool instruction files that never stayed in sync.

But prose has a failure mode this round documented in detail: nothing executes it, so nothing catches its drift. The package shipped a vendored copy of the seven Agent Skills, and every one had drifted from its canonical source - up to 183 changed lines - with one skill still teaching `just` and `pre-commit`, the very stack its own release had removed. I deleted the vendored copy, pointed `AGENTS.md` at the canonical repository, and added a contribution rule with teeth: every command, path, and version named in a skill must exist in the reference implementation.

The boundary matters as much as the feature: `AGENTS.md` guides, it does not enforce. Enforcement stays in types, tests, linters, and branch protection. The file's job is to make an agent's first attempt land close, so the gates have less to catch.

## Trade-offs I Would Not Hide From a Client

- **`ty` is pre-1.0.** It does not yet model the dynamic behavior of Pydantic, Pandera, and MLflow, so I ignore six rule categories globally rather than pretend the code is wrong. It still gates real errors, but mypy remains the fallback.
- **pandas 3.0 is out, and I still cannot use it.** MLflow declares `pandas<3`, and numba caps numpy below 2.5. "Latest everything" is a fantasy; "latest that resolves together" is the job.
- **The gate is slower and the config footprint grew.** `mise run all` lands between one and a half and three and a half minutes on my machine, across nine configuration files the course has to explain. Each one earns its place - but nobody should claim a security scan is free.
- **Four repos in lockstep is still manual discipline.** I found this round's drift by auditing on purpose, a month after swearing everything was aligned. A release tag is a shared checkpoint, not proof of alignment.

## Where This Leaves You

AI agents give us something ML projects have never had: contributors cheap enough to keep a stack current continuously, instead of in painful rewrites every two years. But that only works when the repository is ready for them - one named gate, machine-readable rules, and checks strict enough to catch a confident wrong answer.

Three rules carry most of the value if you are weighing the same jump. Name the gate once, and let hooks, CI, and agents call the same task. Lint everything that executes, and schedule what is slow or flaky - a [lychee](https://lychee.cli.rs) link check over the course found 63 failures of which only 9 were real, a ratio that belongs on a weekly cron, not in a commit gate. And distinguish overriding from suppressing: prove your exceptions with tests, in both directions.

If you are starting a new ML project, generate one from the [cookiecutter](https://github.com/fmind/cookiecutter-mlops-package) and read the [course](https://mlops-coding-course.fmind.dev) chapter behind each choice. If you already have one, the diff from "MLflow 2 + mypy + just" to "MLflow 3 + ty + mise" is smaller than you fear - I just walked the whole path twice, and the commits are public.
