package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func precheckEmailProServiceOrSkip(t *testing.T) string {
	svc := os.Getenv("OVH_EMAILPRO_SERVICE_TEST")
	if svc == "" {
		t.Skip("OVH_EMAILPRO_SERVICE_TEST must be set to run Email Pro acceptance tests (e.g. emailpro-123456)")
	}
	return svc
}

func TestAccEmailProAccount_basic(t *testing.T) {
	service := precheckEmailProServiceOrSkip(t)
	name := "tfacc" + randSuffix(5)
	res := "ovh_email_pro_account.test"

	cfg := fmt.Sprintf(`
provider "ovh" {}

resource "ovh_email_pro_account" "test" {
  service_name = %q
  login        = %q
  password     = "Sup3rSecure!123456"
  display_name = "TF Acc"
  first_name   = "TF"
  last_name    = "Acc"
  quota        = 10240
}
`, service, name)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(res, "service_name", service),
					resource.TestCheckResourceAttr(res, "login", name),
					resource.TestCheckResourceAttr(res, "display_name", "TF Acc"),
					resource.TestCheckResourceAttr(res, "quota", "10240"),
					resource.TestCheckResourceAttrSet(res, "primary_email"),
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
