package openai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectServiceAccountDataSource(t *testing.T) {
	rName := acctest.RandomWithPrefix("openai_tf_test_")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProjectServiceAccountDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.openai_project_service_account.test", "id"),
				),
			},
		},
	})
}

func testAccProjectServiceAccountDataSourceConfig(rName string) string {
	return fmt.Sprintf(`
resource openai_project test {
	name = %[1]q
}
		
resource openai_project_service_account test {
	name = %[2]q
	project_id = openai_project.test.id
}

data "openai_project_service_account" "test" {
	id = openai_project_service_account.test.id
	project_id = openai_project_service_account.test.project_id
}
`, rName, rName)
}
