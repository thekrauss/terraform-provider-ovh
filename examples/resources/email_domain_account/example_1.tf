terraform {
  required_providers {
    ovh = {
      source = "ovh/ovh"
    }
  }
}

provider "ovh" {
  # Credentials are read from environment variables:
  # OVH_ENDPOINT, OVH_APPLICATION_KEY, OVH_APPLICATION_SECRET, OVH_CONSUMER_KEY
}

resource "ovh_email_domain_account" "support" {
  domain       = "mydomain.com"
  account_name = "support"
  description  = "Support mailbox"
  password     = "SecurePassword"
  size         = 5000000000 # ~5 GiB
}
