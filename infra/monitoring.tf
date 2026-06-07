# Observability dashboards, alert policies, notification channels, and uptime checks

# Notification channel for error alerts
resource "google_monitoring_notification_channel" "email" {
  display_name = "Médéric Hurier (Google Chat / Email)"
  type         = "email"
  labels = {
    email_address = var.notification_email
  }

  depends_on = [time_sleep.wait_for_apis]
}


# Alert policy for Cloud Run HTTP 5xx errors
resource "google_monitoring_alert_policy" "cloudrun_5xx_errors" {
  display_name = "Cloud Run 5xx Errors - ${var.service_name}"
  combiner     = "OR"
  conditions {
    display_name = "HTTP 5xx Error Rate > 0"
    condition_threshold {
      filter          = "resource.type = \"cloud_run_revision\" AND resource.labels.service_name = \"${var.service_name}\" AND metric.type = \"run.googleapis.com/request_count\" AND metric.labels.response_code_class = \"5xx\""
      duration        = "60s"
      comparison      = "COMPARISON_GT"
      threshold_value = 0
      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "ALIGN_RATE"
      }
    }
  }

  notification_channels = [
    google_monitoring_notification_channel.email.name
  ]

  alert_strategy {
    auto_close = "172800s" # 2 days
  }
}

# Alert policy for application-level error logs (severity ERROR or higher)
resource "google_monitoring_alert_policy" "error_logs" {
  display_name = "Cloud Run Error Logs - ${var.service_name}"
  combiner     = "OR"
  conditions {
    display_name = "Log matches error condition"
    condition_matched_log {
      filter = "resource.type=\"cloud_run_revision\" AND resource.labels.service_name=\"${var.service_name}\" AND severity>=ERROR"
    }
  }

  notification_channels = [
    google_monitoring_notification_channel.email.name
  ]

  alert_strategy {
    notification_rate_limit {
      period = "300s" # Rate limit to once per 5 minutes to prevent spam
    }
    auto_close = "172800s"
  }
}


