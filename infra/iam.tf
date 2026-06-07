# Service accounts and project-level IAM bindings for run-time identity

# Service Account for Cloud Run
resource "google_service_account" "cloudrun_sa" {
  account_id   = "${var.service_name}-run"
  display_name = "Cloud Run Service Account for ${var.service_name}"
  depends_on   = [time_sleep.wait_for_apis]
}

# IAM permissions for logging
resource "google_project_iam_member" "log_writer" {
  project = var.project_id
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.cloudrun_sa.email}"
}

# IAM permissions for tracing
resource "google_project_iam_member" "trace_agent" {
  project = var.project_id
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.cloudrun_sa.email}"
}


