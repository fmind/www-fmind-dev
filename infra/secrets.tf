# Secret Manager configurations and IAM permissions for application secrets

# Secret for session signing key
resource "google_secret_manager_secret" "session_secret" {
  secret_id = "session-secret"

  replication {
    auto {}
  }

  depends_on = [time_sleep.wait_for_apis]
}

# Grant Cloud Run Service Account access to retrieve the secret value
resource "google_secret_manager_secret_iam_member" "session_secret_accessor" {
  secret_id = google_secret_manager_secret.session_secret.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloudrun_sa.email}"
}
