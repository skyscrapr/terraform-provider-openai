package openai

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFineTuneResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFineTuneResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_finetune.test", "id"),
					resource.TestCheckResourceAttr("openai_finetune.test", "filename", "test.txt"),
					resource.TestCheckResourceAttr("openai_finetune.test", "purpose", "fine-tune"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "openai_finetune.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				// ImportStateVerifyIgnore: []string{"configurable_attribute", "defaulted"},
			},
			// // Update and Read testing
			// {
			// 	Config: testAccExampleResourceConfig("two"),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("scaffolding_example.test", "configurable_attribute", "two"),
			// 	),
			// },
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccFineTuneResourceConfig() string {
	return `
resource "openai_finetune" "example" {
	training_file                  = "file-m5YlZT81Z3kuehmidYGXeo1P"
	validation_file                = "file-n3XhcMU0nyyEphsupzlwOxNx"
	model                          = "ada"
	compute_classification_metrics = true
	classification_positive_class  = " baseball"
}
`
}
