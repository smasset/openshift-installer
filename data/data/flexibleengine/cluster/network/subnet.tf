resource "flexibleengine_vpc_subnet_v1" "subnet" {
  count = var.subnet == null ? 1 : 0

  name          = var.subnet_name
  cidr          = var.subnet_cidr
  gateway_ip    = var.subnet_gateway_ip
  vpc_id        = local.vpc_id
}
