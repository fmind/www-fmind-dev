# Keyless CI: Workload Identity Federation for GitHub Actions.
#
# GitHub's OIDC token impersonates a dedicated service account (no long-lived
# keys) that may push to Artifact Registry and deploy the Cloud Run service.

resource "google_service_account" "github_actions_sa" {
  account_id   = "github-actions"
  display_name = "GitHub Actions Deployment Service Account"
  depends_on   = [time_sleep.wait_for_apis]
}

# Push images to Artifact Registry.
resource "google_artifact_registry_repository_iam_member" "repo_writer" {
  location   = google_artifact_registry_repository.repo.location
  repository = google_artifact_registry_repository.repo.name
  role       = "roles/artifactregistry.writer"
  member     = "serviceAccount:${google_service_account.github_actions_sa.email}"
}

# Deploy new Cloud Run revisions.
resource "google_project_iam_member" "github_actions_run_developer" {
  project = var.project_id
  role    = "roles/run.developer"
  member  = "serviceAccount:${google_service_account.github_actions_sa.email}"
}

# Act as the Cloud Run runtime service account when deploying.
resource "google_service_account_iam_member" "github_actions_act_as_run" {
  service_account_id = google_service_account.cloudrun_sa.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.github_actions_sa.email}"
}

# Workload Identity Pool for GitHub Actions.
resource "google_iam_workload_identity_pool" "github_pool" {
  workload_identity_pool_id = "github-actions-pool"
  display_name              = "GitHub Actions Pool"
  description               = "Identity pool for GitHub Actions CI/CD"
  depends_on                = [time_sleep.wait_for_apis]
}

# Workload Identity Pool Provider for GitHub OIDC.
resource "google_iam_workload_identity_pool_provider" "github_provider" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.github_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "github-provider"
  display_name                       = "GitHub OIDC Provider"
  description                        = "OIDC identity provider for GitHub Actions"

  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.actor"      = "assertion.actor"
    "attribute.ref"        = "assertion.ref"
    "attribute.repository" = "assertion.repository"
  }

  # Mint tokens only for this exact repository (not just any repo in the org), so
  # the pool itself is the first least-privilege gate — not only the downstream
  # principalSet impersonation binding below. assertion.repository is "owner/repo".
  attribute_condition = "assertion.repository == '${var.github_repository}' && assertion.ref == 'refs/heads/main'"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

# Allow only this repository's workflows to impersonate the deployer SA.
resource "google_service_account_iam_member" "wif_impersonate" {
  service_account_id = google_service_account.github_actions_sa.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github_pool.name}/attribute.repository/${var.github_repository}"
}
