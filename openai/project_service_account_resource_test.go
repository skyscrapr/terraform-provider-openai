package openai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectServiceAccountResource_simple(t *testing.T) {
	rName := acctest.RandomWithPrefix("openai_tf_test_")
	resourceName := "openai_project_service_account.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectServiceAccountResourceConfig_simple(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			// ImportState testing
			// {
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	// This is not normally necessary, but is here because this
			// 	// example code does not have an actual upstream service.
			// 	// Once the Read method is able to refresh information from
			// 	// the upstream service, this can be removed.
			// 	// ImportStateVerifyIgnore: []string{"wait"},
			// },
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccProjectServiceAccountResourceConfig_simple(rName string) string {
	return fmt.Sprintf(`
resource openai_project test {
	name = %[1]q
}
		
resource openai_project_service_account test {
	name = %[2]q
	project_id = openai_project.test.id
}
`, rName, rName)
}
