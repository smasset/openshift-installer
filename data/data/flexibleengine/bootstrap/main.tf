provider "flexibleengine" {
  domain_name = var.fe_domain_name
  region      = var.fe_region
  tenant_name = var.fe_project_name
  access_key  = var.fe_access_key
  secret_key  = var.fe_secret_key
  auth_url    = var.fe_auth_url
}
