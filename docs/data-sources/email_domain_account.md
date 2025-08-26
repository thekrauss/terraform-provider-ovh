---
subcategory: "Email"
---



# ovh_email_domain_account (Data Source)

Retrieve information about an existing OVHcloud Web Cloud **email account**.

## Example Usage

```terraform
data "ovh_email_domain_account" "admin" {
  domain       = "mydomain.com"
  account_name = "admin"
}

output "admin_email" {
  value = data.ovh_email_domain_account.admin.email
}
```

## Argument Reference

The following arguments are supported:

* `domain` — (Required) The domain name (for example, `mydomain.com`).
* `account_name` — (Required) Local part of the email address (for example, `support` ⇒ `support@mydomain.com`).

## Attributes Reference

The following attributes are exported:

* `id` — Composite identifier in the form `{domain}/{account_name}`.
* `email` — Full email address (for example, `support@mydomain.com`).
* `description` — The mailbox description.
* `size` — Mailbox size in bytes.
* `is_blocked` — Whether the mailbox is currently blocked.
