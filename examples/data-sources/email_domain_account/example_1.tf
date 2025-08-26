data "ovh_email_domain_account" "admin" {
  domain       = "mydomain.com"
  account_name = "admin"
}

output "admin_email" {
  value = data.ovh_email_domain_account.admin.email
}


