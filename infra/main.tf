# ==============================================================================
# www-fmind-dev Infrastructure Configurations
# ==============================================================================
# The infrastructure configuration has been split into functional files:
#
# - apis.tf        : GCP service APIs enabling and propagation timers
# - registry.tf    : Artifact Registry repositories
# - iam.tf         : Service accounts and general project-level IAM roles
# - cloud_run.tf   : Cloud Run service declarations and service-specific IAM
# - monitoring.tf  : Cloud Monitoring notification channels, dashboards, and alerts
# - providers.tf   : Terraform provider settings and remote backend configuration
# - variables.tf   : Input variables and default settings
# - outputs.tf     : Terraform outputs (e.g. service URLs)
#
# ==============================================================================
