package ovh

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func randSuffix(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestAccEmailDomainAccount_basic(t *testing.T) {
	domain := os.Getenv("OVH_DOMAIN_TEST")
	if domain == "" {
		t.Skip("OVH_DOMAIN_TEST must be set (e.g. mydomain.com)")
	}

	precheckEmailDomainOfferOrSkip(t, domain)

	name := "tfacc" + randSuffix(5)
	res := "ovh_email_domain_account.test"

	cfg := fmt.Sprintf(`
provider "ovh" {}

resource "ovh_email_domain_account" "test" {
  domain       = %q
  account_name = %q
  description  = "tf acceptance"
  password     = "Sup3rSecure!123456"
  size         = 5000000000
}
`, domain, name)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(res, "domain", domain),
					resource.TestCheckResourceAttr(res, "account_name", name),
					resource.TestCheckResourceAttr(res, "size", "5000000000"),
				),
			},
			{
				ResourceName:      res,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	req := []string{
		"OVH_ENDPOINT",
		"OVH_APPLICATION_KEY",
		"OVH_APPLICATION_SECRET",
		"OVH_CONSUMER_KEY",
	}
	for _, k := range req {
		if v := os.Getenv(k); v == "" {
			t.Skipf("%s must be set for acceptance tests", k)
		}
	}
}

func precheckEmailDomainOfferOrSkip(t *testing.T, domain string) {
	client, err := sharedClientForRegion(os.Getenv("OVH_ENDPOINT"))
	if err != nil {
		t.Fatalf("cannot init OVH client: %v", err)
	}
	//the /email/domain/{domain} API returns the service offer.
	var info map[string]any
	endpoint := fmt.Sprintf("/email/domain/%s", url.PathEscape(domain))
	if err := client.Get(endpoint, &info); err != nil {
		t.Logf("warn: cannot fetch %s: %v (will try test anyway)", endpoint, err)
		return
	}
	get := func(k string) string {
		if v, ok := info[k]; ok && v != nil {
			return strings.ToUpper(fmt.Sprint(v))
		}
		return ""
	}
	offer := get("offer")
	if offer == "" {
		offer = get("type")
	} // fallback
	if offer == "REDIRECT" {
		t.Skipf("Domain %s has Redirect offer; cannot create mailboxes. Use a domain with MX Plan/Email Pro.", domain)
	}
}
