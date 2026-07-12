# Terraform required providers and GCS state backend settings.
#
terraform {
  required_version = ">= 1.1.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 7.35"
    }
    time = {
      source  = "hashicorp/time"
      version = "~> 0.14"
    }
  }

  backend "gcs" {
    bucket = "www-fmind-dev-tfstate"
    prefix = "infra/state"
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}
