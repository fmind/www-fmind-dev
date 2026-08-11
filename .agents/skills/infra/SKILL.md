---
name: infra
description: Change the Google Cloud infrastructure for www.fmind.dev with OpenTofu — validate offline, plan, and apply manually against live state. Use for any change under infra/.
license: MIT
metadata:
  author: Médéric HURIER (Fmind)
---

# Change the Infrastructure

`infra/` is a flat OpenTofu root module owning every cloud resource behind <https://www.fmind.dev/>: the Cloud Run service, Artifact Registry, the runtime and CI service accounts, Workload Identity Federation, error alerting, and the cookieless BigQuery analytics route. State lives in the versioned GCS bucket `www-fmind-dev-tfstate` under prefix `infra/state`.

**Applying is always a deliberate human step.** No task and no hook runs `tofu apply` — plans move real infrastructure and real money.

## Division of Ownership

CI and OpenTofu both touch the Cloud Run service, and they do not fight:

- **CI owns the image.** `.github/workflows/deploy.yml` deploys an immutable digest on every `main` push.
- **OpenTofu owns the shape** — CPU, memory, scaling, env vars, probes, IAM. The image tag is under `lifecycle.ignore_changes`, so a plan never rolls the container back to whatever digest the state remembers.

An infrastructure change therefore does **not** trigger a deployment, and a deployment does **not** drift the infrastructure.

## Workflow

1. **Validate, without touching state**:
   ```bash
   mise run format          # tofu fmt -recursive, among the rest
   mise run check           # includes check:scan
   mise run check:tofu      # tofu fmt -check, init -backend=false, validate, tflint
   ```
   `check:scan` runs `trivy config` over the module and the Dockerfile; it is offline and part of every commit. `check:tofu` is not: `init` downloads provider schemas, and in a working copy where step 4 has already run a real `init`, it also opens the gcs backend and needs live credentials. That is why it is not in the pre-commit hook — its own `infra` workflow gates every `infra/` change on a fresh checkout, where neither applies.

1. **Authenticate** for anything that reads real state:
   ```bash
   gcloud auth application-default login
   ```

1. **Plan and read it** — never skip reading:
   ```bash
   tofu -chdir=infra init
   tofu -chdir=infra plan -out=tmp/plan.tfplan
   ```

1. **Apply the reviewed plan**, deliberately:
   ```bash
   tofu -chdir=infra apply tmp/plan.tfplan
   ```

## Gotchas

1. **First OpenTofu init migrates the lockfile**: the module moved from HashiCorp Terraform to OpenTofu, so provider source addresses are now `registry.opentofu.org/...`. The first `tofu init` against the existing GCS state reconciles that. The resources are unchanged — a plan right after the migration must come back empty. **If it does not, stop and read it before applying.**
1. **State is secret**: it stores every attribute in plaintext. It never belongs in git — `.gitignore` blocks `*.tfstate*` and `*.tfvars` — only in the versioned bucket.
1. **Provider majors move fast**: `google` and `google-beta` are pinned exactly (`= 7.43.0`) because resources rename across majors. Bump them deliberately with the [upgrade-tools](~/.agents/skills/upgrade-tools/SKILL.md) skill and read the upgrade guide.
1. **The privacy boundary is code, not config**: the external edge must strip inbound `X-Client-Geo` before injecting trusted geography. `middleware.go` treats that edge rule as the boundary — changing the edge means changing both.
1. **Keyless only**: CI authenticates through branch-restricted Workload Identity Federation. Never create a service-account key.
