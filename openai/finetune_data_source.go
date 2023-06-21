package openai

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ModelDataSource{}

func NewFineTuneDataSource() datasource.DataSource {
	return &FineTuneDataSource{OpenAIDatasource: &OpenAIDatasource{}}
}

// FineTineDataSource defines the data source implementation.
type FineTuneDataSource struct {
	*OpenAIDatasource
}

func (d *FineTuneDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_finetune"
}

func (d *FineTuneDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fine-Tine data source",

		Attributes: openAIFineTuneAttributes(),
	}
}

func (d *FineTuneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OpenAIFineTuneModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	fineTune, err := d.client.FineTunes().GetFineTune(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read FineTune, got error: %s", err))
		return
	}

	data = NewOpenAIFineTuneModel(fineTune)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func openAIFineTuneAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "File Identifier",
			Required:            true,
		},
		"object": schema.StringAttribute{
			MarkdownDescription: "Object Type",
			Computed:            true,
		},
		"model": schema.StringAttribute{
			MarkdownDescription: "Model ID",
			Computed:            true,
		},
		"created": schema.Int64Attribute{
			MarkdownDescription: "Created Time",
			Computed:            true,
		},
		"events": schema.ListNestedAttribute{
			MarkdownDescription: "Events",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"created": schema.Int64Attribute{
						MarkdownDescription: "Created Time",
						Computed:            true,
					},
					"object": schema.StringAttribute{
						MarkdownDescription: "Object Type",
						Computed:            true,
					},
					"level": schema.StringAttribute{
						MarkdownDescription: "Level",
						Computed:            true,
					},
					"message": schema.StringAttribute{
						MarkdownDescription: "Message",
						Computed:            true,
					},
				},
			},
		},
		"hyperparams": schema.SingleNestedAttribute{
			MarkdownDescription: "Hyperparams",
			Computed:            true,
			Attributes: map[string]schema.Attribute{
				"batch_size": schema.Int64Attribute{
					MarkdownDescription: "Batch Size",
					Computed:            true,
				},
				"learning_rate_multiplier": schema.Float64Attribute{
					MarkdownDescription: "Learning Rate Multipier",
					Computed:            true,
				},
				"n_epochs": schema.Int64Attribute{
					MarkdownDescription: "N Epochs",
					Computed:            true,
				},
				"prompt_loss_weight": schema.Float64Attribute{
					MarkdownDescription: "Prompt Loss Weight",
					Computed:            true,
				},
			},
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
		"result_files": schema.ListNestedAttribute{
			MarkdownDescription: "Result Files",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: openAIFileDatasourceAttributes(),
			},
		},
		"validation_files": schema.ListNestedAttribute{
			MarkdownDescription: "Validation Files",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: openAIFileDatasourceAttributes(),
			},
		},
		"training_files": schema.ListNestedAttribute{
			MarkdownDescription: "Training Files",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: openAIFileDatasourceAttributes(),
			},
		},
		"updated_at": schema.Int64Attribute{
			MarkdownDescription: "Updated Time",
			Computed:            true,
		},
	}
}

func openAIFileDatasourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "File Identifier",
			Required:            true,
		},
		"bytes": schema.Int64Attribute{
			MarkdownDescription: "File size in bytes",
			Computed:            true,
		},
		"created": schema.Int64Attribute{
			MarkdownDescription: "Created Time",
			Computed:            true,
		},
		"filename": schema.StringAttribute{
			MarkdownDescription: "Filename",
			Computed:            true,
		},
		"filepath": schema.StringAttribute{
			MarkdownDescription: "Filepath",
			Computed:            true,
		},
		"object": schema.StringAttribute{
			MarkdownDescription: "Object Type",
			Computed:            true,
		},
		"purpose": schema.StringAttribute{
			MarkdownDescription: "Intended use of file. Use 'fine-tune' for Fine-tuning",
			Computed:            true,
		},
	}
}
