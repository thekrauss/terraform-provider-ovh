package ovh

import (
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEmailDomainAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEmailDomainAccountsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {Type: schema.TypeString, Computed: true},
						"email":        {Type: schema.TypeString, Computed: true},
						"description":  {Type: schema.TypeString, Computed: true},
						"size":         {Type: schema.TypeInt, Computed: true},
						"is_blocked":   {Type: schema.TypeBool, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceEmailDomainAccountsRead(d *schema.ResourceData, meta any) error {
	config := meta.(*Config)
	client := config.OVHClient

	domain := d.Get("domain").(string)

	//lists account names
	var accountNames []string
	endpoint := fmt.Sprintf("/email/domain/%s/account", url.PathEscape(domain))
	if err := client.Get(endpoint, &accountNames); err != nil {
		return err
	}

	//recovers the details of each account
	accounts := make([]map[string]any, 0, len(accountNames))
	for _, accountName := range accountNames {
		var resp emailAccountResp
		endpoint := fmt.Sprintf("/email/domain/%s/account/%s", url.PathEscape(domain), url.PathEscape(accountName))
		if err := client.Get(endpoint, &resp); err != nil {
			return err
		}
		accounts = append(accounts, map[string]any{
			"account_name": resp.AccountName,
			"email":        resp.Email,
			"description":  resp.Description,
			"size":         resp.Size,
			"is_blocked":   resp.IsBlocked,
		})
	}

	//sets the values
	d.SetId(domain) //the ID of the source data = domain name
	_ = d.Set("accounts", accounts)

	return nil
}
