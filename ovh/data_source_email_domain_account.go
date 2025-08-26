package ovh

import (
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Reads information from an existing mailbox without creating/modifying it.
func dataSourceEmailDomainAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEmailDomainAccountRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_blocked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceEmailDomainAccountRead(d *schema.ResourceData, meta any) error {
	config := meta.(*Config)
	client := config.OVHClient

	domain := d.Get("domain").(string)
	accountName := d.Get("account_name").(string)

	var resp emailAccountResp
	endpoint := fmt.Sprintf("/email/domain/%s/account/%s", url.PathEscape(domain), url.PathEscape(accountName))
	if err := client.Get(endpoint, &resp); err != nil {
		return err
	}

	//the Terraform ID will be "domain/account_name"
	d.SetId(fmt.Sprintf("%s/%s", domain, accountName))
	_ = d.Set("description", resp.Description)
	_ = d.Set("size", resp.Size)
	_ = d.Set("email", resp.Email)
	_ = d.Set("is_blocked", resp.IsBlocked)

	return nil
}
