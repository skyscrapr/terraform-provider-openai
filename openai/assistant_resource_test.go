package openai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAssistantResource_tool_code_interpreter(t *testing.T) {
	rName := acctest.RandomWithPrefix("openai_tf_test_")
	assistantResourceName := "openai_assistant.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAssistantResourceConfig_tool_code_interpreter("./test-fixtures/test.jsonl", rName, "test description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(assistantResourceName, "id"),
					resource.TestCheckResourceAttr(assistantResourceName, "name", rName),
					resource.TestCheckResourceAttr(assistantResourceName, "description", "test description"),
					resource.TestCheckResourceAttr(assistantResourceName, "model", "gpt-4"),
					resource.TestCheckResourceAttr(assistantResourceName, "instructions", "You are a personal math tutor. When asked a question, write and run Python code to answer the question."),
				),
			},
			// Update and Read testing
			{
				Config: testAccAssistantResourceConfig_tool_code_interpreter("./test-fixtures/test.jsonl", rName, "test description updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(assistantResourceName, "description", "test description updated"),
				),
			},
			// ImportState testing
			{
				ResourceName:      assistantResourceName,
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

func TestAccAssistantResource_tool_function(t *testing.T) {
	rName := acctest.RandomWithPrefix("openai_tf_test_")
	assistantResourceName := "openai_assistant.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAssistantResourceConfig_tool_function("./test-fixtures/test.jsonl", rName, "test description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(assistantResourceName, "id"),
					resource.TestCheckResourceAttr(assistantResourceName, "name", rName),
					resource.TestCheckResourceAttr(assistantResourceName, "description", "test description"),
					resource.TestCheckResourceAttr(assistantResourceName, "model", "gpt-3.5-turbo-0125"),
					resource.TestCheckResourceAttr(assistantResourceName, "instructions", "You are the personal assistant for users who are using our app."),
				),
			},
			// Update and Read testing
			{
				Config: testAccAssistantResourceConfig_tool_function("./test-fixtures/test.jsonl", rName, "test description updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(assistantResourceName, "description", "test description updated"),
				),
			},
			// ImportState testing
			{
				ResourceName:      assistantResourceName,
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

func testAccAssistantResourceConfig_tool_code_interpreter(filename string, rName string, description string) string {
	return fmt.Sprintf(`	
resource openai_file test {
	filepath = %[1]q
}

resource openai_assistant test {
	name = %[2]q
	description = %[3]q
	model = "gpt-4"
	instructions = "You are a personal math tutor. When asked a question, write and run Python code to answer the question."
	tools = [
		{type = "code_interpreter"}
	]
	tool_resources = {
		code_interpreter = {
			file_ids = [
				openai_file.test.id,
			]
		}
	}
}
`, filename, rName, description)
}

func testAccAssistantResourceConfig_tool_function(filename string, rName string, description string) string {
	return fmt.Sprintf(`
resource openai_assistant test {
	name = %[2]q
	description = %[3]q
	model = "gpt-3.5-turbo-0125"
	instructions = "You are the personal assistant for users who are using our app."
	tools = [
		{
			type = "function"
			function = {
				name = "get_additional_info"
				description = "Get additional information."
				parameters = jsonencode({
				  type = "object",
				  properties = {}
				})
			}
		}
	]
}
`, filename, rName, description)
}
