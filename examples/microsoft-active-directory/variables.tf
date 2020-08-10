// Provider setup

variable "onefuse_address" {
  type = string
}

variable "onefuse_port" {
  type = string
  default = "443"
}

variable "onefuse_user" {
  type = string
}

variable "onefuse_password" {
  type = string
}

variable "onefuse_verify_ssl" {
  type = bool
  default = false
}

// Microsoft AD Endpoint

variable "onefuse_microsoft_endpoint" {
  type = string
  default = "My Microsoft Endpoint"
}


// Microsoft AD Policy

variable "ad_policy_name" {
  type = string
  default = "Some Naming Policy"
}

variable "ad_policy_description" {
  type = string
  default = "Created with Terraform"
}

variable "ad_computer_name_letter_case" {
  type = string
  default = "Lowercase"
}

variable "ad_ou" {
  type = string
  default = "ou=Accounting,dc=onefuse,dc=com"
}

