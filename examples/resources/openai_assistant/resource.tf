terraform {
  required_providers {
    openai = {
      source = "skyscrapr/openai"
    }
  }
}

resource "openai_assistant" "test" {
  name         = "tf-test-sample"
  description  = "tf-test-sample"
  model        = "gpt-4"
  instructions = "You are a personal math tutor. When asked a question, write and run Python code to answer the question."
  tools = [
    {
      type = "code_interpreter"
    }
  ]
}