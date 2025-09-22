package ovh

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

type emailAccountResp struct {
	AccountName string `json:"accountName"`
	Description string `json:"description"`
	Domain      string `json:"domain"`
	Email       string `json:"email"`
	IsBlocked   bool   `json:"isBlocked"`
	Size        int    `json:"size"`
}

func resourceEmailDomainAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceEmailDomainAccountCreate,
		Read:   resourceEmailDomainAccountRead,
		Delete: resourceEmailDomainAccountDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5000000000,
				ForceNew: true,
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

func resourceEmailDomainAccountCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.OVHClient
	domain := d.Get("domain").(string)
	body := map[string]any{
		"accountName": d.Get("account_name").(string),
		"password":    d.Get("password").(string),
		"size":        d.Get("size").(int),
	}
	if v, ok := d.GetOk("description"); ok {
		body["description"] = v.(string)
	}

	endpoint := fmt.Sprintf("/email/domain/%s/account", url.PathEscape(domain))
	if err := client.Post(endpoint, body, nil); err != nil {
		if apiErr, ok := err.(*ovh.APIError); ok && apiErr.Code == 403 && strings.Contains(strings.ToLower(err.Error()), "redirect") {
			return fmt.Errorf("ovh_email_domain_account: domain %s has a Redirect email offer; a MX Plan is required to create accounts", domain)
		}
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", domain, d.Get("account_name").(string)))
	return resourceEmailDomainAccountRead(d, meta)
}

func resourceEmailDomainAccountRead(d *schema.ResourceData, meta any) error {
	config := meta.(*Config)
	client := config.OVHClient

	domain, accountName, err := parseEmailAccountID(d.Id())
	if err != nil {
		return err
	}

	var resp emailAccountResp
	endpoint := fmt.Sprintf("/email/domain/%s/account/%s", domain, url.PathEscape(accountName))
	if err := client.Get(endpoint, &resp); err != nil {
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	_ = d.Set("domain", domain)
	_ = d.Set("account_name", resp.AccountName)
	_ = d.Set("description", resp.Description)
	_ = d.Set("size", resp.Size)
	_ = d.Set("email", resp.Email)
	_ = d.Set("is_blocked", resp.IsBlocked)

	return nil
}

func resourceEmailDomainAccountDelete(d *schema.ResourceData, meta any) error {
	config := meta.(*Config)
	client := config.OVHClient

	domain, accountName, err := parseEmailAccountID(d.Id())
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("/email/domain/%s/account/%s", domain, url.PathEscape(accountName))
	if err := client.Delete(endpoint, nil); err != nil {

		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId("")
	return nil
}

func parseEmailAccountID(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid ID %q (expected format: domain/account_name)", id)
	}
	return parts[0], parts[1], nil
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*ovh.APIError); ok {
		return apiErr.Code == http.StatusNotFound
	}
	return strings.Contains(strings.ToLower(err.Error()), "404")
}
