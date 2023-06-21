package openai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFileResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFileResourceConfig("./test-fixtures/test.jsonl"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_file.test", "id"),
					resource.TestCheckResourceAttr("openai_file.test", "filename", "test.jsonl"),
					resource.TestCheckResourceAttr("openai_file.test", "purpose", "fine-tune"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "openai_file.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"filepath"},
			},
			// // Update and Read testing
			// {
			// 	Config: testAccFileResourceConfig("two"),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("scaffolding_example.test", "configurable_attribute", "two"),
			// 	),
			// },
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccFileResourceConfig(filename string) string {
	return fmt.Sprintf(`	
resource "openai_file" "test" {
	filepath = %[1]q
}
`, filename)
}
