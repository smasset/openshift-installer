resource "flexibleengine_vpc_v1" "vpc" {
  count = var.vpc == null ? 1 : 0

  name        = var.vpc_name
  cidr        = var.vpc_cidr
  description = local.description
}
