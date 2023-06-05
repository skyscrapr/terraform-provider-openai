terraform {
  required_providers {
    openai = {
      source = "skyscrapr/openai"
    }
  }
}

provider "openai" {}

data "openai_files" "test" {
}

data "openai_file" "test" {
  id = "1"
}

