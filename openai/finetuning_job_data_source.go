package openai

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ModelDataSource{}

func NewFineTuningJobDataSource() datasource.DataSource {
	return &FineTuningJobDataSource{OpenAIDatasource: &OpenAIDatasource{}}
}

// FineTuningJobDataSource defines the data source implementation.
type FineTuningJobDataSource struct {
	*OpenAIDatasource
}

func (d *FineTuningJobDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_finetuning_job"
}

func (d *FineTuningJobDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fine-Tuning Job data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Fine Tuning Job Identifier",
				Required:            true,
			},
			"object": schema.StringAttribute{
				MarkdownDescription: "Object Type",
				Computed:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "Created Time",
				Computed:            true,
			},
			"finished_at": schema.Int64Attribute{
				MarkdownDescription: "Finished Time",
				Computed:            true,
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "Model ID",
				Computed:            true,
			},
			"fine_tuned_model": schema.StringAttribute{
				MarkdownDescription: "Fine-Tuned Model ID",
				Computed:            true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Status",
				Computed:            true,
			},
			"hyperparams": schema.SingleNestedAttribute{
				MarkdownDescription: "Hyperparams",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"n_epochs": schema.Int64Attribute{
						MarkdownDescription: "N Epochs",
						Computed:            true,
					},
				},
			},
			"validation_file": schema.StringAttribute{
				MarkdownDescription: "Validation File",
				Computed:            true,
			},
			"training_file": schema.StringAttribute{
				MarkdownDescription: "Training File",
				Computed:            true,
			},
			"result_files": schema.ListAttribute{
				MarkdownDescription: "Result Files",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"trained_tokens": schema.Int64Attribute{
				MarkdownDescription: "Trained Tokens",
				Computed:            true,
			},
		},
	}
}

func (d *FineTuningJobDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OpenAIFineTuningJobModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	fineTune, err := d.client.FineTuning().GetFineTuningJob(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read FineTune, got error: %s", err))
		return
	}

	data = NewOpenAIFineTuningJobModel(fineTune)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
