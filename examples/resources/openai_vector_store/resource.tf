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

resource "openai_vector_store" "test" {
  name = "test_vector_store"
  file_ids = [
    openai_file.test.id
  ]
}
