terraform {
  required_providers {
    openai = {
      source = "skyscrapr/openai"
    }
  }
}

# resource "openai_assistant" "test" {
#   name         = "tf-test-sample"
#   description  = "tf-test-sample"
#   model        = "gpt-4"
#   instructions = "You are a personal math tutor. When asked a question, write and run Python code to answer the question."
#   tools = [
#     {
#       type = "code_interpreter"
#     }
#   ]
# }

resource "openai_file" "test" {
  filepath = "./test-fixtures/test.jsonl"
}

resource "openai_vector_store" "test" {
  name = "test_vector_store"
  file_ids = [
    openai_file.test.id
  ]
}

resource "openai_assistant" "test" {
  name         = "test_assistant"
  description  = "my test assistant description"
  model        = "gpt-3.5-turbo"
  instructions = "You are a personal math tutor. When asked a question, write and run Python code to answer the question."
  tools = [
    { type = "file_search" }
  ]
  tool_resources = {
    file_search = {
      vector_store_ids = [
        openai_vector_store.test.id,
      ]
    }
  }
}