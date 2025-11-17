package ovh

import (
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEmailProAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEmailProAccountRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"login": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"primary_email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"quota": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"task_pending_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceEmailProAccountRead(d *schema.ResourceData, meta any) error {
	config := meta.(*Config)
	client := config.OVHClient

	service := d.Get("service_name").(string)
	login := d.Get("login").(string)

	var account emailProAccount

	endpoint := fmt.Sprintf("/email/pro/%s/account/%s", url.PathEscape(service), url.PathEscape(login))
	if err := client.Get(endpoint, &account); err != nil {
		if isNotFound(err) {
			return fmt.Errorf("account %s/%s not found", service, login)
		}
		return explainEmailProErr(service, login, err)

	}

	d.SetId(fmt.Sprintf("%s/%s", service, login))

	_ = d.Set("primary_email", account.PrimaryEmailAddress)
	_ = d.Set("display_name", account.DisplayName)
	_ = d.Set("first_name", account.FirstName)
	_ = d.Set("last_name", account.LastName)
	_ = d.Set("quota", account.Quota)
	_ = d.Set("state", account.State)
	_ = d.Set("task_pending_id", account.TaskPendingId)

	return nil
}
