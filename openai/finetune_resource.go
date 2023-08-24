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
var _ resource.Resource = &FineTuneResource{}
var _ resource.ResourceWithImportState = &FineTuneResource{}

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
			"wait": schema.BoolAttribute{
				MarkdownDescription: "Wait for Fine-Tune completion",
				Optional:            true,
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
			"learning_rate_multiplier": schema.Float64Attribute{
				MarkdownDescription: "Learning Rate Multiplier",
				Optional:            true,
			},
			"prompt_loss_weight": schema.Float64Attribute{
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
			"fine_tune": schema.SingleNestedAttribute{
				MarkdownDescription: "FineTune",
				Computed:            true,
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Fine Tune Identifier",
						Computed:            true,
					},
					"object": schema.StringAttribute{
						MarkdownDescription: "Object Type",
						Computed:            true,
					},
					"model": schema.StringAttribute{
						MarkdownDescription: "Model Identifier",
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
					"hyperparams": schema.SingleNestedAttribute{
						MarkdownDescription: "Hyperparams",
						Computed:            true,
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
					"organization_id": schema.StringAttribute{
						MarkdownDescription: "Organization Id",
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
							Attributes: openAIFileResourceAttributes(),
						},
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
			},
		},
	}
}

func (r *FineTuneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OpenAIFineTuneResourceModel
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
		LearningRateMultiplier:       data.LearningRateMultiplier.ValueFloat64(),
		PromptLossWeight:             data.PromptLossWeight.ValueFloat64(),
		ComputeClassificationMetrics: data.ComputeClassificationMetrics.ValueBool(),
		ClassificationNClasses:       data.ClassificationNClasses.ValueInt64(),
		ClassificationPositiveClass:  data.ClassificationPositiveClass.ValueString(),
		ClassificationBetas:          data.ClassificationBetas,
		Suffix:                       data.Suffix,
	}
	fineTune, err := r.client.FineTunes().CreateFineTune(&ftreq)
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to create fine tuning job, got error: %s", err))
		return
	}
	tflog.Info(ctx, "FineTune created successfully")
	data.Id = types.StringValue(fineTune.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if !data.Wait.IsUnknown() && data.Wait.ValueBool() {
		tflog.Info(ctx, "Begin Streaming")
		err = r.client.FineTunes().SubscribeFineTuneEvents(
			fineTune.Id,
			func(event *openai.FineTuneEvent) error {
				tflog.Info(ctx, fmt.Sprintf("Fine -Tune Event: %s", event.Message))
				return nil
			},
		)
		if err != nil {
			resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to stream fine tuning events, got error: %s", err))
			return
		}
	}

	// Update finetune prior to saving to state
	fineTune, err = r.client.FineTunes().GetFineTune(fineTune.Id)
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read fine tuning job, got error: %s", err))
		return
	}

	data.FineTune, _ = types.ObjectValueFrom(ctx, data.FineTune.AttributeTypes(ctx), NewOpenAIFineTuneModel(fineTune))
	data.Id = types.StringValue(fineTune.Id)
	if len(fineTune.TrainingFiles) > 0 {
		data.TrainingFile = types.StringValue(fineTune.TrainingFiles[0].Id)
	}
	if len(fineTune.ValidationFiles) > 0 {
		data.ValidationFile = types.StringValue(fineTune.ValidationFiles[0].Id)
	}
	data.Model = types.StringValue(fineTune.Model)
	// data.NEpochs = types.Int64Value(fineTune.Hyperparams.NEpochs)
	// data.BatchSize = types.Int64Value(fineTune.Hyperparams.BatchSize)
	// data.LearningRateMultiplier = types.Float64Value(fineTune.Hyperparams.LearningRateMultiplier)
	// data.PromptLossWeight = types.Float64Value(fineTune.Hyperparams.PromptLossWeight)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FineTuneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OpenAIFineTuneResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Reading Fine-Tune with id: %s", data.Id.ValueString()))
	fineTune, err := r.client.FineTunes().GetFineTune(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read Fine Tune, got error: %s", err))
		return
	}

	data.FineTune, _ = types.ObjectValueFrom(ctx, data.FineTune.AttributeTypes(ctx), NewOpenAIFineTuneModel(fineTune))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FineTuneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "Update not supported.")
}

func (r *FineTuneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OpenAIFineTuneResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Get existing Fine-Tune...")
	fineTune, err := r.client.FineTunes().GetFineTune(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read Fine Tune, got error: %s", err))
		return
	}
	// need to ignore is fine tune is already delete. - Missing test

	// Cancel fine tune
	tflog.Info(ctx, fmt.Sprintf("Fine-Tune.Status: %s", fineTune.Status))
	switch fineTune.Status {
	case "succeeded", "cancelled", "failed":
	default:
		tflog.Info(ctx, "Cancelling Fine-Tune")
		_, err = r.client.FineTunes().CancelFineTune(data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to cancel fine tune %s, got error: %s", fineTune.Id, err))
			return
		}
	}

	// Delete result files
	for _, file := range fineTune.ResultFiles {
		tflog.Info(ctx, fmt.Sprintf("Deleting Fine-Tune Result File: %s", file.Id))
		_, err := r.client.Files().DeleteFile(file.Id)
		if err != nil {
			apiError := GetOpenAIAPIError(err)
			if apiError != nil && apiError.HTTPStatusCode == 404 {
				tflog.Info(ctx, "Fine-Tune Result File does not exist")
				err = nil
			}
		}
		if err != nil {
			resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to delete Result File %s, got error: %s", file.Id, err))
			return
		}
	}

	// Delete the fine tuned model
	if fineTune.FineTunedModel != "" {
		tflog.Info(ctx, fmt.Sprintf("Deleting Fine-Tune Model: %s", fineTune.FineTunedModel))
		bDeleted, err := r.client.FineTunes().DeleteFineTuneModel(fineTune.FineTunedModel)
		if err != nil {
			if err, ok := err.(*openai.APIError); ok {
				fmt.Println("openai error:", err.Code)
				// Or whatever other field(s) you need
			}

			resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to delete Fime Tune Model, got error: %s", err))
			return
		}
		if !bDeleted {
			tflog.Trace(ctx, "Fine Tune Model not deleted")
		}
		tflog.Trace(ctx, "Fine Tune Model deleted successfully")
	}
}

func (r *FineTuneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
