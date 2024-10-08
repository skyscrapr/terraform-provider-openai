---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "openai_file Resource - terraform-provider-openai"
subcategory: ""
description: |-
  File resource
---

# openai_file (Resource)

File resource

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `filepath` (String) Filename

### Optional

- `purpose` (String) Intended use of file. Use 'fine-tune' for Fine-tuning

### Read-Only

- `bytes` (Number) File size in bytes
- `created` (Number) Created Time
- `filename` (String) Filename
- `id` (String) File Identifier
- `object` (String) Object Type
