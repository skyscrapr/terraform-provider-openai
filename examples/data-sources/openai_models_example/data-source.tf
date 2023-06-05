terraform {
  required_providers {
    openai = {
      source = "skyscrapr/openai"
    }
  }
}

provider "openai" {}

data "openai_models" "test" {
}

data "openai_model" "test" {
  id = "whisper-1"
}

