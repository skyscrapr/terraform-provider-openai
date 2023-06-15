terraform {
  required_providers {
    openai = {
      source = "skyscrapr/openai"
    }
  }
}

provider "openai" {}

data "openai_finetune" "example" {
  id = "ft-UV2XKz7N4T9O5WkB7ojzmg6a"
}
