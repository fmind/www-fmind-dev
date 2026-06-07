# Configuration for Cloud Run container service deployment and IAM access policies

# Cloud Run Service
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

      startup_probe {
        initial_delay_seconds = 0
        timeout_seconds       = 3
        period_seconds        = 5
        failure_threshold     = 3
        http_get {
          path = "/healthz"
        }
      }

    }
  }

  lifecycle {
    ignore_changes = [
      client,
      client_version,
      template[0].containers[0].image,
    ]
  }

  depends_on = [
    time_sleep.wait_for_apis,
    google_project_iam_member.log_writer,
    google_project_iam_member.trace_agent,
  ]
}

# Public access
resource "google_cloud_run_v2_service_iam_member" "noauth" {
  location = var.region
  name     = google_cloud_run_v2_service.web.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
