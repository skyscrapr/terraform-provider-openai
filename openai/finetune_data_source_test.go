package openai

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFineTuneDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccFineTuneDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.openai_finetune.test", "id"),
				),
			},
		},
	})
}

const testAccFineTuneDataSourceConfig = `
data "openai_finetune" "test" {
	id = "ft-UV2XKz7N4T9O5WkB7ojzmg6a"
}
`
