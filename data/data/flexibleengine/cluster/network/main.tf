locals {
  vpc_id      = var.vpc != null ? data.flexibleengine_vpc_v1.vpc[0].id : flexibleengine_vpc_v1.vpc[0].id
  subnet_id   = var.subnet != null ? data.flexibleengine_vpc_subnet_v1.subnet[0].id : flexibleengine_vpc_subnet_v1.subnet[0].id
  description = "Created By OpenShift Installer"
}
