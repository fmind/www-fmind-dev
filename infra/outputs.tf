# Terraform outputs for resource access and verification
#
output "github_actions_service_account_email" {
  description = "The service account email for GitHub Actions"
  value       = google_service_account.github_actions_sa.email
}

output "repository_url" {
  description = "The URL of the Artifact Registry repository"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${var.repository_id}"
}

output "service_url" {
  description = "The URL of the deployed Cloud Run service"
  value       = google_cloud_run_v2_service.web.uri
}

output "wif_provider_name" {
  description = "The Workload Identity Provider name for GitHub Actions"
  value       = google_iam_workload_identity_pool_provider.github_provider.name
}

