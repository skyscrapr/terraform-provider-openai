package openai

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProjectDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.openai_project.test", "id"),
				),
			},
		},
	})
}

const testAccProjectDataSourceConfig = `
data "openai_projects" "test" {}

data "openai_project" "test" {
	id = data.openai_projects.test.projects[0].id
}
`
