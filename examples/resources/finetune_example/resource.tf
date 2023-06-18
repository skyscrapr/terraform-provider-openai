terraform {
  required_providers {
    openai = {
      source = "skyscrapr/openai"
    }
  }
}

provider "openai" {}

resource "openai_file" "training_file" {
  filename = "sport2_prepared_train.jsonl"
}

resource "openai_file" "validation_file" {
  filename = "sport2_prepared_valid.jsonl"
}

resource "openai_finetune" "example" {
  training_file                  = openai_file.training_file.id
  validation_file                = openai_file.validation_file.id
  model                          = "ada"
  compute_classification_metrics = true
  classification_positive_class  = " baseball"
  wait                           = true
}

