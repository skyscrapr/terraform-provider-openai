package openai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFileDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccFileDataSourceConfig("./test-fixtures/test.jsonl"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.openai_file.test", "id"),
					resource.TestCheckResourceAttr("data.openai_file.test", "filename", "test.jsonl"),
				),
			},
		},
	})
}

func testAccFileDataSourceConfig(filename string) string {
	return fmt.Sprintf(`	
resource "openai_file" "test" {
	filepath = %[1]q
}

data "openai_file" "test" {
  id = openai_file.test.id
}
`, filename)
}
