---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "openai_files Data Source - terraform-provider-openai"
subcategory: ""
description: |-
  Files data source
---

# openai_files (Data Source)

Files data source



<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `files` (Attributes List) Files (see [below for nested schema](#nestedatt--files))
- `id` (String) Files identifier

<a id="nestedatt--files"></a>
### Nested Schema for `files`

Required:

- `id` (String) File Identifier

Read-Only:

- `bytes` (Number) File size in bytes
- `created` (Number) Created Time
- `filename` (String) Filename
- `object` (String) Object Type
- `purpose` (String) Intended use of file. Use 'fine-tune' for Fine-tuning