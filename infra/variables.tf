# Infrastructure variables defining default settings and parameters.

variable "github_repository" {
  description = "The GitHub repository in the format owner/repo"
  type        = string
  default     = "fmind/www-fmind-dev"
}

variable "google_analytics_id" {
  description = "Google Analytics 4 measurement ID injected into the Cloud Run service (production only). Public value; empty disables analytics."
  type        = string
  default     = "G-Z28QKD99V2"
}

variable "image_uri" {
  description = "The Docker image URI to deploy. Defaults to a hello-world image for the initial create; the live image is managed out-of-band by CI (see cloud_run.tf lifecycle.ignore_changes)."
  type        = string
  default     = "us-docker.pkg.dev/cloudrun/container/hello"
}

variable "notification_email" {
  description = "The email address for alert notifications"
  type        = string
  default     = "mederic.hurier@fmind.dev"
}

variable "project_id" {
  description = "The GCP project ID to deploy resources to"
  type        = string
  default     = "www-fmind-dev"
}

variable "region" {
  description = "The GCP region to deploy resources to"
  type        = string
  default     = "us-central1"
}

variable "repository_id" {
  description = "The name of the Artifact Registry repository"
  type        = string
  default     = "app"
}

variable "service_name" {
  description = "The name of the Cloud Run service"
  type        = string
  default     = "www-fmind-dev"
}
