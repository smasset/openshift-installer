data "flexibleengine_vpc_v1" "vpc" {
  count = var.vpc != null ? 1 : 0

  name = var.vpc
}

data "flexibleengine_vpc_subnet_v1" "subnet" {
  count = var.subnet != null ? 1 : 0

  name = var.subnet
}
