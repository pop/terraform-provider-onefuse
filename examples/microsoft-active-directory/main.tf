provider "onefuse" {
  scheme      = "http"
  address     = var.onefuse_address
  port        = var.onefuse_port
  user        = var.onefuse_user
  password    = var.onefuse_password
  verify_ssl  = var.onefuse_verify_ssl
}

data "onefuse_microsoft_endpoint" "my_microsoft_endpoint" {
    name = var.onefuse_microsoft_endpoint
}

resource "onefuse_microsoft_ad_policy" "my_ad_policy" {
    name = var.ad_policy_name
    description = var.ad_policy_description
    microsoft_endpoint_id = data.onefuse_microsoft_endpoint.my_microsoft_endpoint.id
    computer_name_letter_case = var.ad_computer_name_letter_case
    ou = var.ad_ou
    workspace_url = var.ad_workspace_url
}
