terraform {
  required_providers {
    openai = {
      source = "skyscrapr/openai"
    }
  }
}

provider "openai" {}

data "openai_finetune" "example" {
  id = "ft-JsPjykPvJNSHdkvyOmfNo3hA"
}
