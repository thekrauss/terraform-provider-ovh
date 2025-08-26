package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataEmailDomainAccount_basic(t *testing.T) {
	domain := os.Getenv("OVH_DOMAIN_TEST")
	if domain == "" {
		t.Skip("OVH_DOMAIN_TEST must be set (e.g. mydomain.com)")
	}

	precheckEmailDomainOfferOrSkip(t, domain)

	name := "tfacc" + randSuffix(5)

	cfg := fmt.Sprintf(`
provider "ovh" {}

resource "ovh_email_domain_account" "test" {
  domain       = %q
  account_name = %q
  description  = "tfacc - data source read"
  password     = "Secure!123456"
  size         = 5000000000
}

data "ovh_email_domain_account" "read" {
  domain       = ovh_email_domain_account.test.domain
  account_name = ovh_email_domain_account.test.account_name
}
`, domain, name)

	dataRes := "data.ovh_email_domain_account.read"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					//The data ID = "domain/account"
					resource.TestCheckResourceAttr(dataRes, "id", fmt.Sprintf("%s/%s", domain, name)),
					//The full email returned by the API
					resource.TestCheckResourceAttr(dataRes, "email", fmt.Sprintf("%s@%s", name, domain)),
					//The size must match that of the resource
					resource.TestCheckResourceAttr(dataRes, "size", "5000000000"),
					//Calculated fields must be present
					resource.TestCheckResourceAttrSet(dataRes, "description"),
					resource.TestCheckResourceAttrSet(dataRes, "is_blocked"),
				),
			},
		},
	})
}
