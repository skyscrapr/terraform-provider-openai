terraform {
  required_providers {
    openai = {
      source = "skyscrapr/openai"
    }
  }
}

resource "openai_file" "test" {
  filepath = "./test-fixtures/test.jsonl"
}
