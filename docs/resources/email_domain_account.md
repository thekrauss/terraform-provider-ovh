---
subcategory : "Email"
---



# ovh_email_domain_account

Creates and manages an OVHcloud Web Cloud **email account** for a given domain.

## Example Usage

```terraform
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
```

## Argument Reference

The following arguments are supported:

* `domain` — (Required) The domain name (for example, `mydomain.com`).
* `account_name` — (Required) Local part of the email address (for example, `support` ⇒ `support@mydomain.com`).
* `password` — (Required, Sensitive) Password for the mailbox.
* `size` — (Optional) Mailbox size in **bytes**. Default: `5000000000` (≈ 5 GiB).
* `description` — (Optional) Free-form description.

## Attributes Reference

The following attributes are exported:

* `id` — Composite identifier in the form `{domain}/{account_name}`.
* `email` — Full email address (for example, `support@mydomain.com`).
* `is_blocked` — Whether the mailbox is currently blocked.
* `domain` — The domain name.
* `account_name` — The local part.
* `size` — Mailbox size in bytes.
* `description` — The mailbox description.

## Import

Existing accounts can be imported using the composite identifier:

```bash
terraform import ovh_email_domain_account.support mydomain.com/support
