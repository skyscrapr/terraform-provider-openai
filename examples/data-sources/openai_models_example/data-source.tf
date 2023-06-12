terraform {
  required_providers {
    openai = {
      source = "skyscrapr/openai"
    }
  }
}

data "openai_models" "models" {
}

output "model_count" {
  value = length(data.openai_models.models.models)
}