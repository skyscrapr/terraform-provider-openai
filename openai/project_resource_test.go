package openai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectResource_tool_simple(t *testing.T) {
	rName := acctest.RandomWithPrefix("openai_tf_test_")
	projectResourceName := "openai_project.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfig_simple(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(projectResourceName, "id"),
					resource.TestCheckResourceAttr(projectResourceName, "name", rName),
				),
			},
			// Update and Read testing
			{
				Config: testAccProjectResourceConfig_simple(rName + "_U"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(projectResourceName, "name", rName+"_U"),
				),
			},
			// ImportState testing
			{
				ResourceName:      projectResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				// ImportStateVerifyIgnore: []string{"wait"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccProjectResourceConfig_simple(rName string) string {
	return fmt.Sprintf(`	
resource openai_project test {
	name = %[1]q
}
`, rName)
}
