---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "openai_project Data Source - terraform-provider-openai"
subcategory: ""
description: |-
  Project data source
---

# openai_project (Data Source)

Project data source



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The identifier, which can be referenced in API endpoints

### Read-Only

- `archived_at` (Number) The Unix timestamp (in seconds) of when the project was archived or null.
- `created_at` (Number) The Unix timestamp (in seconds) of when the project was created.
- `name` (String) The name of the project. This appears in reporting.
- `object` (String) The object type, which is always organization.project
- `status` (String) active or archived
