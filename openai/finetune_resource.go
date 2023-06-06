package openai

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/skyscrapr/openai-sdk-go/openai"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FileResource{}
var _ resource.ResourceWithImportState = &FileResource{}

func NewFineTuneResource() resource.Resource {
	return &FineTuneResource{OpenAIResource: &OpenAIResource{}}
}

// FineTuneResource defines the resource implementation.
type FineTuneResource struct {
	*OpenAIResource
}

func (r *FineTuneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_finetune"
}

func (r *FineTuneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fine Tune resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Fine Tune Identifier",
				Computed:            true,
			},
			"training_file": schema.StringAttribute{
				MarkdownDescription: "Training File Identifier",
				Optional:            true,
			},
			"validation_file": schema.StringAttribute{
				MarkdownDescription: "Validation File Identifier",
				Optional:            true,
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "Model Identifier",
				Optional:            true,
			},
			"n_epochs": schema.Int64Attribute{
				MarkdownDescription: "N Epochs",
				Optional:            true,
			},
			"batch_size": schema.Int64Attribute{
				MarkdownDescription: "Batch Size",
				Optional:            true,
			},
			"learning_rate_multiplier": schema.Int64Attribute{
				MarkdownDescription: "Learning Rate Multiplier",
				Optional:            true,
			},
			"prompt_loss_weight": schema.Int64Attribute{
				MarkdownDescription: "Prompt Loss Weight",
				Optional:            true,
			},
			"compute_classification_metrics": schema.BoolAttribute{
				MarkdownDescription: "Compute Classification Metrics",
				Optional:            true,
			},
			"classification_n_classes": schema.Int64Attribute{
				MarkdownDescription: "Classification N Classes",
				Optional:            true,
			},
			"classification_positive_class": schema.StringAttribute{
				MarkdownDescription: "Classification Positive Class",
				Optional:            true,
			},
			"classification_betas": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Classification Betas",
				Optional:            true,
			},
			"suffix": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Suffix",
				Optional:            true,
			},
			"object": schema.StringAttribute{
				MarkdownDescription: "Object Type",
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
						"object": schema.StringAttribute{
							MarkdownDescription: "Object Type",
							Computed:            true,
						},
						"created": schema.Int64Attribute{
							MarkdownDescription: "Created Time",
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
			"fine_tuned_model": schema.StringAttribute{
				MarkdownDescription: "Fine Tuned Model",
				Computed:            true,
			},
			"hyperparams": schema.ListNestedAttribute{
				MarkdownDescription: "Hyperparams",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"batch_size": schema.Int64Attribute{
							MarkdownDescription: "Batch Size",
							Computed:            true,
						},
						"learning_rate_multiplier": schema.Float64Attribute{
							MarkdownDescription: "Learning Rate Multiplier",
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
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization Id",
				Computed:            true,
			},
			"result_files": schema.ListNestedAttribute{
				MarkdownDescription: "Result Files",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: openAIFileResourceAttributes(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Status",
				Computed:            true,
			},
			"validation_files": schema.ListNestedAttribute{
				MarkdownDescription: "Validation Files",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: openAIFileResourceAttributes(),
				},
			},
			"training_files": schema.ListNestedAttribute{
				MarkdownDescription: "Training Files",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: openAIFileResourceAttributes(),
				},
			},
			"updated_at": schema.Int64Attribute{
				MarkdownDescription: "Updated Time",
				Computed:            true,
			},
		},
	}
}

type OpenAICreateFineTuneRequest struct {
	TrainingFile                 types.String `tfsdk:"training_file"`
	ValidationFile               types.String `tfsdk:"training_file"`
	Model                        types.String `tfsdk:"model"`
	NEpochs                      types.Int64  `tfsdk:"n_epochs"`
	BatchSize                    types.Int64  `tfsdk:"batch_size"`
	LearningRateMultiplier       types.Int64  `tfsdk:"learning_rate_multiplier"`
	PromptLossWeight             types.Int64  `tfsdk:"prompt_loss_weight"`
	ComputeClassificationMetrics types.Bool   `tfsdk:"compute_classification_metrics"`
	ClassificationNClasses       types.Int64  `tfsdk:"classification_n_classes"`
	ClassificationPositiveClass  types.String `tfsdk:"classification_positive_class"`
	ClassificationBetas          []string     `tfsdk:"classification_betas"`
	Suffix                       []string     `tfsdk:"suffix"`
}

func (r *FineTuneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OpenAICreateFineTuneRequest

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ftreq := openai.CreateFineTunesRequest{
		TrainingFile:                 data.TrainingFile.ValueString(),
		ValidationFile:               data.ValidationFile.ValueString(),
		Model:                        data.Model.ValueString(),
		NEpochs:                      data.NEpochs.ValueInt64(),
		BatchSize:                    data.BatchSize.ValueInt64(),
		LearningRateMultiplier:       data.LearningRateMultiplier.ValueInt64(),
		PromptLossWeight:             data.PromptLossWeight.ValueInt64(),
		ComputeClassificationMetrics: data.ComputeClassificationMetrics.ValueBool(),
		ClassificationNClasses:       data.ClassificationNClasses.ValueInt64(),
		ClassificationPositiveClass:  data.ClassificationPositiveClass.ValueString(),
		ClassificationBetas:          data.ClassificationBetas,
		Suffix:                       data.Suffix,
	}
	fineTune, err := r.client.FineTunes().CreateFineTune(&ftreq)
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to upload& File, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "Uploaded file successfully")

	fineTuneData := NewOpenAIFineTuneModel(fineTune)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &fineTuneData)...)
}

func (r *FineTuneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OpenAIFineTuneModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	fineTune, err := r.client.FineTunes().GetFineTune(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read Fine Tune, got error: %s", err))
		return
	}

	data = NewOpenAIFineTuneModel(fineTune)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FineTuneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "Update not supported.")
}

func (r *FineTuneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OpenAIFineTuneModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bDeleted, err := r.client.FineTunes().DeleteFineTuneModel(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to delete Fime Tune Model, got error: %s", err))
		return
	}
	if bDeleted {
		tflog.Trace(ctx, "Fine Tune Model deleted successfully")
	} else {
		tflog.Trace(ctx, "Fine Tune Model not deleted")
	}
}

func (r *FineTuneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func openAIFileResourceAttributes() map[string]schema.Attribute {
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
