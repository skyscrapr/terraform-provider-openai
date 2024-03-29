package openai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFineTuningJobResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccOpenAI(t); testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFineTuneResourceConfig("./test-fixtures/test_prepared_train.jsonl", "./test-fixtures/test_prepared_valid.jsonl"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_finetuning_job.test", "id"),
					resource.TestCheckResourceAttrSet("openai_finetuning_job.test", "training_file"),
					resource.TestCheckResourceAttrSet("openai_finetuning_job.test", "validation_file"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "openai_finetuning_job.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"wait"},
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

func testAccFineTuneResourceConfig(training_file string, validation_file string) string {
	return fmt.Sprintf(`	
resource openai_file training_file {
	filepath = %[1]q
}

resource openai_file validation_file {
	filepath = %[2]q
}

resource "openai_finetuning_job" "test" {
	training_file                  = openai_file.training_file.id
	validation_file                = openai_file.validation_file.id
	model                          = "babbage-002"
	wait = true
}
`, training_file, validation_file)
}
