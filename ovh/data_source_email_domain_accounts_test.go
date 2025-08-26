package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataEmailDomainAccounts_list(t *testing.T) {
	domain := os.Getenv("OVH_DOMAIN_TEST")
	if domain == "" {
		t.Skip("OVH_DOMAIN_TEST must be set (e.g. mydomain.com)")
	}
	precheckEmailDomainOfferOrSkip(t, domain)

	nameA := "tfaccA" + randSuffix(5)
	nameB := "tfaccB" + randSuffix(5)

	cfg := fmt.Sprintf(`
provider "ovh" {}

# On crée deux boîtes
resource "ovh_email_domain_account" "a" {
  domain       = %q
  account_name = %q
  description  = "tfacc plural A"
  password     = "Sup3rSecure!123456"
  size         = 5000000000
}

resource "ovh_email_domain_account" "b" {
  domain       = %q
  account_name = %q
  description  = "tfacc plural B"
  password     = "Sup3rSecure!123456"
  size         = 5000000000
}

# On liste toutes les boîtes du domaine
data "ovh_email_domain_accounts" "all" {
  domain = %q
}
`, domain, nameA, domain, nameB, domain)

	dataRes := "data.ovh_email_domain_accounts.all"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					//The data ID = the domain
					resource.TestCheckResourceAttr(dataRes, "id", domain),
					//Checks that the two accounts created are present in the list (order indifferent)
					checkAccountsContain(dataRes, domain, []string{nameA, nameB}),
				),
			},
		},
	})
}

// checkAccountsContain verifies that the "accounts" source data contains
// the expected accounts (by their account_name) and that the emails match.
func checkAccountsContain(resName, domain string, expectedNames []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms, ok := s.RootModule().Resources[resName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resName)
		}
		attrs := ms.Primary.Attributes

		//Number of items in the list
		nStr, ok := attrs["accounts.#"]
		if !ok {
			return fmt.Errorf("attribute not found: accounts.#")
		}
		n, err := atoi(nStr)
		if err != nil {
			return fmt.Errorf("invalid accounts.#: %v", err)
		}

		//Index the account_name present
		present := map[string]struct{}{}
		for i := 0; i < n; i++ {
			pfx := fmt.Sprintf("accounts.%d.", i)
			aname := attrs[pfx+"account_name"]
			if aname != "" {
				present[aname] = struct{}{}
			}
		}

		//Check presence and consistent email for each expected
		for _, want := range expectedNames {
			if _, ok := present[want]; !ok {
				return fmt.Errorf("expected account %q not found in data source", want)
			}
			//Look for the index where account_name == want to check the email
			foundEmail := false
			for i := 0; i < n; i++ {
				pfx := fmt.Sprintf("accounts.%d.", i)
				if attrs[pfx+"account_name"] == want {
					email := attrs[pfx+"email"]
					if email != fmt.Sprintf("%s@%s", want, domain) {
						return fmt.Errorf("email mismatch for %q: got %q, want %q",
							want, email, fmt.Sprintf("%s@%s", want, domain))
					}
					foundEmail = true
					break
				}
			}
			if !foundEmail {
				return fmt.Errorf("email not found for account %q", want)
			}
		}
		return nil
	}
}

func atoi(s string) (int, error) {
	var x int
	_, err := fmt.Sscan(s, &x)
	return x, err
}
