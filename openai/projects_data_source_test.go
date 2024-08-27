package openai

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProjectsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify ID has any value set
					resource.TestCheckResourceAttrSet("data.openai_projects.test", "id"),
				),
			},
		},
	})
}

const testAccProjectsDataSourceConfig = `
data "openai_projects" "test" {}
`
