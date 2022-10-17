variable "vpc" {
  type        = string
  description = "The pre-existing FE VPC where the cluster will be created."
}

variable "vpc_name" {
  type        = string
  description = "The new FE VPC name where the cluster will be created."
}

variable "vpc_cidr" {
  type        = string
  description = "The new FE VPC CIDR where the cluster will be created."
}

variable "subnet" {
  type        = string
  description = "The pre-existing FE Subnet where the cluster will be created."
}

variable "subnet_name" {
  type        = string
  description = "The new FE Subnet name where the cluster will be created."
}

variable "subnet_cidr" {
  type        = string
  description = "The new FE Subnet CIDR where the cluster will be created."
}

variable "subnet_gateway_ip" {
  type        = string
  description = "The new FE Subnet gateway IP where the cluster will be created."
}
