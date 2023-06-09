---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "openai_finetune Data Source - terraform-provider-openai"
subcategory: ""
description: |-
  Fine-Tine data source
---

# openai_finetune (Data Source)

Fine-Tine data source



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) File Identifier

### Read-Only

- `created` (Number) Created Time
- `events` (Attributes List) Events (see [below for nested schema](#nestedatt--events))
- `fine_tuned_model` (String) Fine-Tuned Model ID
- `hyperparams` (Attributes) Hyperparams (see [below for nested schema](#nestedatt--hyperparams))
- `model` (String) Model ID
- `object` (String) Object Type
- `organization_id` (String) Organization ID
- `result_files` (Attributes List) Result Files (see [below for nested schema](#nestedatt--result_files))
- `status` (String) Status
- `training_files` (Attributes List) Training Files (see [below for nested schema](#nestedatt--training_files))
- `updated_at` (Number) Updated Time
- `validation_files` (Attributes List) Validation Files (see [below for nested schema](#nestedatt--validation_files))

<a id="nestedatt--events"></a>
### Nested Schema for `events`

Read-Only:

- `created` (Number) Created Time
- `level` (String) Level
- `message` (String) Message
- `object` (String) Object Type


<a id="nestedatt--hyperparams"></a>
### Nested Schema for `hyperparams`

Read-Only:

- `batch_size` (Number) Batch Size
- `learning_rate_multiplier` (Number) Learning Rate Multipier
- `n_epochs` (Number) N Epochs
- `prompt_loss_weight` (Number) Prompt Loss Weight


<a id="nestedatt--result_files"></a>
### Nested Schema for `result_files`

Required:

- `id` (String) File Identifier

Read-Only:

- `bytes` (Number) File size in bytes
- `created` (Number) Created Time
- `filename` (String) Filename
- `filepath` (String) Filepath
- `object` (String) Object Type
- `purpose` (String) Intended use of file. Use 'fine-tune' for Fine-tuning


<a id="nestedatt--training_files"></a>
### Nested Schema for `training_files`

Required:

- `id` (String) File Identifier

Read-Only:

- `bytes` (Number) File size in bytes
- `created` (Number) Created Time
- `filename` (String) Filename
- `filepath` (String) Filepath
- `object` (String) Object Type
- `purpose` (String) Intended use of file. Use 'fine-tune' for Fine-tuning


<a id="nestedatt--validation_files"></a>
### Nested Schema for `validation_files`

Required:

- `id` (String) File Identifier

Read-Only:

- `bytes` (Number) File size in bytes
- `created` (Number) Created Time
- `filename` (String) Filename
- `filepath` (String) Filepath
- `object` (String) Object Type
- `purpose` (String) Intended use of file. Use 'fine-tune' for Fine-tuning
