package openai

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFineTuneDataSource(t *testing.T) {
	t.Skip("TODO")
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
	id = "1"
}
`
