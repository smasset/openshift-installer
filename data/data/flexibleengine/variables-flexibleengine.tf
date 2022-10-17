variable "fe_domain_name" {
  type        = string
  description = "The target FE domain for the cluster."
}

variable "fe_region" {
  type        = string
  description = "The target FE region for the cluster."
}

variable "fe_project_name" {
  type        = string
  description = "The target FE project for the cluster."
  default     = null
}

variable "fe_access_key" {
  type        = string
  description = "The FE access key to interact with the FE APIs."
}

variable "fe_secret_key" {
  type        = string
  description = "The FE secret key to interact with the FE APIs."
}

variable "fe_auth_url" {
  type        = string
  description = "The URL used to authenticate to interact with the FE APIs."
  default     = "https://iam.eu-west-0.prod-cloud-ocb.orange-business.com/v3"
}

variable "fe_vpc" {
  type        = string
  description = "The pre-existing FE VPC where the cluster will be created."
  default     = null
}

variable "fe_vpc_name" {
  type        = string
  description = "The new FE VPC name where the cluster will be created."
  default     = null
}

variable "fe_vpc_cidr" {
  type        = string
  description = "The new FE VPC CIDR where the cluster will be created."
  default     = null
}

variable "fe_subnet" {
  type        = string
  description = "The pre-existing FE Subnet where the cluster will be created."
  default     = null
}

variable "fe_subnet_name" {
  type        = string
  description = "The new FE Subnet name where the cluster will be created."
  default     = null
}

variable "fe_subnet_cidr" {
  type        = string
  description = "The new FE Subnet CIDR where the cluster will be created."
  default     = null
}

variable "fe_subnet_gateway_ip" {
  type        = string
  description = "The new FE Subnet gateway IP where the cluster will be created."
  default     = null
}
