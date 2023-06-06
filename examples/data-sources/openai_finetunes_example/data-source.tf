terraform {
  required_providers {
    openai = {
      source = "skyscrapr/openai"
    }
  }
}

provider "openai" {}

data "openai_finetunes" "test" {

}
