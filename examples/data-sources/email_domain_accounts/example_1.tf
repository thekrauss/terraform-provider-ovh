#  all mailboxes for a domain
data "ovh_email_domain_accounts" "all" {
  domain = "mydomain.com"
}

# extract only the emails
output "all_emails" {
  value = [for a in data.ovh_email_domain_accounts.all.accounts : a.email]
}

# show first account (if any)
output "first_account" {
  value = length(data.ovh_email_domain_accounts.all.accounts) > 0
    ? data.ovh_email_domain_accounts.all.accounts[0]
    : null
}
