# Cloud Run service deployment and its public-access IAM policy.

resource "google_cloud_run_v2_service" "web" {
  name                = var.service_name
  location            = var.region
  ingress             = "INGRESS_TRAFFIC_ALL"
  deletion_protection = true

  template {
    service_account = google_service_account.cloudrun_sa.email

    # No request legitimately runs long (the app's own WriteTimeout is 10s), so
    # cap Cloud Run's request timeout well below the 300s default to fail fast.
    timeout = "30s"

    scaling {
      # Scale to zero when idle (min 0); allow a few instances so a traffic
      # spike is absorbed instead of throttled against a single-container ceiling.
      max_instance_count = 3
      min_instance_count = 0
    }

    containers {
      image = var.image_uri
      ports {
        container_port = 8080
      }
      resources {
        limits = {
          # 256Mi is ample for this static Go server (small binary + embedded
          # assets, no heavy runtime); leaves comfortable GC headroom under burst.
          cpu    = "1"
          memory = "256Mi"
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
      startup_probe {
        initial_delay_seconds = 0
        timeout_seconds       = 3
        period_seconds        = 5
        failure_threshold     = 3
        http_get {
          path = "/health"
        }
      }
      liveness_probe {
        timeout_seconds   = 3
        period_seconds    = 30
        failure_threshold = 3
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
