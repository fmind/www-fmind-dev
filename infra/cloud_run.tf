# Cloud Run service deployment and its public-access IAM policy.

resource "google_cloud_run_v2_service" "web" {
  name     = var.service_name
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = google_service_account.cloudrun_sa.email

    scaling {
      max_instance_count = 1
      min_instance_count = 0
    }

    containers {
      image = var.image_uri
      ports {
        container_port = 8080
      }
      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
        cpu_idle          = true
        startup_cpu_boost = true
      }

      # Runtime configuration for the Go binary (Twelve-Factor env). The image
      # also defaults ENVIRONMENT=production, but declaring it here keeps the
      # deployed contract explicit and drift-visible.
      env {
        name  = "ENVIRONMENT"
        value = "production"
      }
      env {
        name  = "GOOGLE_ANALYTICS_ID"
        value = var.google_analytics_id
      }

      startup_probe {
        initial_delay_seconds = 0
        timeout_seconds       = 3
        period_seconds        = 5
        failure_threshold     = 3
        http_get {
          path = "/health"
        }
      }
    }
  }

  lifecycle {
    ignore_changes = [
      client,
      client_version,
      # The live image is rolled by CI (build -> push -> deploy); Terraform owns
      # the service shape, not the image tag.
      template[0].containers[0].image,
      template[0].labels,
    ]
  }

  depends_on = [
    time_sleep.wait_for_apis,
    google_project_iam_member.log_writer,
    google_project_iam_member.trace_agent,
  ]
}

# Public, unauthenticated access (this is a public website).
resource "google_cloud_run_v2_service_iam_member" "noauth" {
  location = var.region
  name     = google_cloud_run_v2_service.web.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
