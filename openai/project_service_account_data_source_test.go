package openai

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectServiceAccountDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProjectServiceAccountDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.openai_project_service_account.test", "id"),
				),
			},
		},
	})
}

const testAccProjectServiceAccountDataSourceConfig = `
data "openai_project_service_accounts" "test" {}

data "openai_project_service_account" "test" {
	id = data.openai_project_service_accounts.test.project_service_accounts[0].id
}
`
