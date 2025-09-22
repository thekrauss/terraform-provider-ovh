package ovh

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

type emailProAccount struct {
	Login               string `json:"login"`
	PrimaryEmailAddress string `json:"primaryEmailAddress"`
	DisplayName         string `json:"displayName"`
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	Quota               int    `json:"quota"`
	State               string `json:"state"`
	TaskPendingId       int    `json:"taskPendingId"`
}

type emailProAccountCreate struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName,omitempty"`
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Quota       int    `json:"quota,omitempty"`
}

func resourceEmailProAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceEmailProAccountCreate,
		Read:   resourceEmailProAccountRead,
		Delete: resourceEmailProAccountDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
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

			// Computed
			"primary_email": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceEmailProAccountCreate(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*Config)
	c := cfg.OVHClient

	service := d.Get("service_name").(string)

	body := emailProAccountCreate{
		Login:    d.Get("login").(string),
		Password: d.Get("password").(string),
	}
	if v, ok := d.GetOk("display_name"); ok {
		body.DisplayName = v.(string)
	}
	if v, ok := d.GetOk("first_name"); ok {
		body.FirstName = v.(string)
	}
	if v, ok := d.GetOk("last_name"); ok {
		body.LastName = v.(string)
	}
	if v, ok := d.GetOk("quota"); ok {
		body.Quota = v.(int)
	}

	endpoint := fmt.Sprintf("/email/pro/%s/account", url.PathEscape(service))
	if err := c.Post(endpoint, &body, nil); err != nil {
		return explainEmailProErr(service, body.Login, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", service, body.Login))
	return resourceEmailProAccountRead(d, meta)
}

func resourceEmailProAccountRead(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*Config)
	c := cfg.OVHClient

	service, login, err := splitTwo(d.Id())
	if err != nil {
		return err
	}

	var acc emailProAccount
	endpoint := fmt.Sprintf("/email/pro/%s/account/%s", url.PathEscape(service), url.PathEscape(login))
	if err := c.Get(endpoint, &acc); err != nil {
		if isOvhNotFound(err) {
			d.SetId("")
			return nil
		}
		return explainEmailProErr(service, login, err)
	}

	_ = d.Set("service_name", service)
	_ = d.Set("login", acc.Login)
	_ = d.Set("display_name", acc.DisplayName)
	_ = d.Set("first_name", acc.FirstName)
	_ = d.Set("last_name", acc.LastName)
	_ = d.Set("quota", acc.Quota)

	_ = d.Set("primary_email", acc.PrimaryEmailAddress)
	_ = d.Set("state", acc.State)
	_ = d.Set("task_pending_id", acc.TaskPendingId)

	return nil
}

func resourceEmailProAccountDelete(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*Config)
	c := cfg.OVHClient

	service, login, err := splitTwo(d.Id())
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("/email/pro/%s/account/%s", url.PathEscape(service), url.PathEscape(login))
	if err := c.Delete(endpoint, nil); err != nil {
		if isOvhNotFound(err) {
			d.SetId("")
			return nil
		}
		return explainEmailProErr(service, login, err)
	}

	d.SetId("")
	return nil
}

func splitTwo(id string) (string, string, error) {
	p := strings.Split(id, "/")
	if len(p) != 2 || p[0] == "" || p[1] == "" {
		return "", "", fmt.Errorf("invalid id %q, expected service/login", id)
	}
	return p[0], p[1], nil
}

func isOvhNotFound(err error) bool {
	if apiErr, ok := err.(*ovh.APIError); ok {
		return apiErr.Code == 404
	}
	return strings.Contains(strings.ToLower(err.Error()), "404")
}

func explainEmailProErr(service, login string, err error) error {
	if apiErr, ok := err.(*ovh.APIError); ok {
		if apiErr.Code == 403 {
			return fmt.Errorf("email_pro_account(%s/%s): forbidden (403). Make sure the domain is attached to the Email Pro service, that the requested quota is valid, and that the Consumer Key has the rights /email/pro/*: %v", service, login, err)
		}
	}
	return err
}
