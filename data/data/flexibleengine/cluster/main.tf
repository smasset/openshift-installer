provider "flexibleengine" {
  domain_name = var.fe_domain_name
  region      = var.fe_region
  tenant_name = var.fe_project_name
  access_key  = var.fe_access_key
  secret_key  = var.fe_secret_key
  auth_url    = var.fe_auth_url
}

module "network" {
  source = "./network"

  vpc               = var.fe_vpc
  vpc_name          = var.fe_vpc_name
  vpc_cidr          = var.fe_vpc_cidr
  subnet            = var.fe_subnet
  subnet_name       = var.fe_subnet_name
  subnet_cidr       = var.fe_subnet_cidr
  subnet_gateway_ip = var.fe_subnet_gateway_ip
}
