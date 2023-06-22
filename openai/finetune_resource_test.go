package openai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFineTuneResource(t *testing.T) {
	t.Skip("Cost associated with test")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccOpenAI(t); testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFineTuneResourceConfig("./test-fixtures/test_prepared_train.jsonl", "./test-fixtures/test_prepared_valid.jsonl"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_finetune.test", "id"),
					resource.TestCheckResourceAttrSet("openai_finetune.test", "training_file"),
					resource.TestCheckResourceAttrSet("openai_finetune.test", "validation_file"),
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
				ImportStateVerifyIgnore: []string{"compute_classification_metrics", "classification_n_classes", "classification_positive_class", "classification_betas", "suffix"},
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

resource "openai_finetune" "test" {
	training_file                  = openai_file.training_file.id
	validation_file                = openai_file.validation_file.id
	model                          = "ada"
	compute_classification_metrics = true
	classification_positive_class  = " baseball"
	wait = true
}
`, training_file, validation_file)
}
