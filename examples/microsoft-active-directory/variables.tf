// Provider setup

variable "onefuse_scheme" {
  type = string
  default = "https"
}

variable "onefuse_address" {
  type = string
  default = "localhost"
}

variable "onefuse_port" {
  type = string
  default = "443"
}

variable "onefuse_user" {
  type = string
  default = "admin"
}

variable "onefuse_password" {
  type = string
  default = "admin"
}

variable "onefuse_verify_ssl" {
  type = bool
  default = false
}

// Microsoft AD Endpoint

variable "onefuse_microsoft_endpoint" {
  type = string
  default = "myMicrosoftEndpoint"
}

// Microsoft AD Policy

variable "ad_policy_name" {
  type = string
  default = "myADPolicy"
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

variable "ad_workspace_url" {
  type = string
  default = ""
}
