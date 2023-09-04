package openai

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/skyscrapr/openai-sdk-go/openai"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FineTuningJobResource{}
var _ resource.ResourceWithImportState = &FineTuningJobResource{}

func NewFineTuningJobResource() resource.Resource {
	return &FineTuningJobResource{OpenAIResource: &OpenAIResource{}}
}

// FineTuningJobResource defines the resource implementation.
type FineTuningJobResource struct {
	*OpenAIResource
}

func (r *FineTuningJobResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_finetuning_job"
}

func (r *FineTuningJobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fine Tuning Job resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Fine Tuning Job Identifier",
				Computed:            true,
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
				MarkdownDescription: "Model Identifier",
				Optional:            true,
			},
			"fine_tuned_model": schema.StringAttribute{
				MarkdownDescription: "Fine Tuned Model",
				Computed:            true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization Id",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Status",
				Computed:            true,
			},
			"hyperparams": schema.SingleNestedAttribute{
				MarkdownDescription: "Hyperparams",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"n_epochs": schema.StringAttribute{
						MarkdownDescription: "N Epochs",
						Optional:            true,
					},
				},
			},
			"training_file": schema.StringAttribute{
				MarkdownDescription: "Training File Identifier",
				Optional:            true,
			},
			"validation_file": schema.StringAttribute{
				MarkdownDescription: "Validation File Identifier",
				Optional:            true,
			},
			"suffix": schema.StringAttribute{
				MarkdownDescription: "Suffix",
				Optional:            true,
			},
			"result_files": schema.ListNestedAttribute{
				MarkdownDescription: "Result Files",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: openAIFileResourceAttributes(),
				},
			},
			"trained_tokens": schema.Int64Attribute{
				MarkdownDescription: "Trained Tokens",
				Computed:            true,
			},
			"wait": schema.BoolAttribute{
				MarkdownDescription: "Wait for Fine Tuning Job completion",
				Optional:            true,
			},
		},
	}
}

func (r *FineTuningJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OpenAIFineTuningJobModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating FineTuning Job...")

	ftreq := openai.CreateFineTuningJobRequest{
		TrainingFile:   data.TrainingFile.ValueString(),
		ValidationFile: data.ValidationFile.ValueString(),
		Model:          data.Model.ValueString(),
		Suffix:         data.Suffix.ValueString(),
	}
	data.Hyperparams.As(ctx, ftreq.Hyperparameters, basetypes.ObjectAsOptions{})

	createTimeout := 100 * time.Hour

	var ftJob *openai.FineTuningJob
	var err error
	err = retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		ftJob, err = r.client.FineTuning().CreateFineTuningJob(&ftreq)
		if err != nil {
			apiError := GetOpenAIAPIError(err)
			if apiError != nil && apiError.HTTPStatusCode == 400 {
				tflog.Info(ctx, fmt.Sprintf("%s - Retrying...", err))
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to create fine tuning job, got error: %s", err))
		return
	}
	tflog.Info(ctx, "FineTuning Job created successfully")
	data = NewOpenAIFineTuningJobModel(ftJob)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if !data.Wait.IsUnknown() && data.Wait.ValueBool() {
		tflog.Info(ctx, "Waiting for fine tuning job completion...")
		var lastEvent *string
		lastEvent = nil

		err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
			events, err := r.client.FineTuning().ListFineTuningEvents(ftJob.Id, lastEvent, nil)
			if err != nil {
				return retry.NonRetryableError(err)
			}
			tflog.Trace(ctx, "Fine Tuning Events obtained successfully")

			for _, event := range events {
				tflog.Info(ctx, fmt.Sprintf("Fine-Tuning Event: %s", event.Message))
				lastEvent = &event.Id
			}
			// Update finetuning job state
			ftJob, err = r.client.FineTuning().GetFineTuningJob(ftJob.Id)
			if err != nil {
				return retry.NonRetryableError(err)
			}
			data = NewOpenAIFineTuningJobModel(ftJob)
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

			switch ftJob.Status {
			case "succeeded":
				return nil
			case "running":
				return retry.RetryableError(fmt.Errorf("fine tuning job still running"))
			default:
				return retry.NonRetryableError(fmt.Errorf("unexpected job status: %s", ftJob.Status))
			}
		})
		if err != nil {
			resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to stream fine tuning events, got error: %s", err))
			return
		}
	}
}

func (r *FineTuningJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OpenAIFineTuningJobModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Reading Fine-Tune with id: %s", data.Id.ValueString()))
	ftJob, err := r.client.FineTuning().GetFineTuningJob(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read Fine Tuning Job, got error: %s", err))
		return
	}

	data = NewOpenAIFineTuningJobModel(ftJob)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FineTuningJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "Update not supported.")
}

func (r *FineTuningJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OpenAIFineTuningJobModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Get existing Fine-Tune...")
	ftJob, err := r.client.FineTuning().GetFineTuningJob(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read Fine Tune, got error: %s", err))
		return
	}
	// need to ignore is fine tune is already delete. - Missing test

	// Cancel fine tune
	tflog.Info(ctx, fmt.Sprintf("Fine-Tuning-Job.Status: %s", ftJob.Status))
	switch ftJob.Status {
	case "succeeded", "cancelled", "failed":
	default:
		tflog.Info(ctx, "Cancelling Fine-Tune")
		_, err = r.client.FineTuning().CancelFineTuningJob(data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to cancel fine tune %s, got error: %s", ftJob.Id, err))
			return
		}
	}

	// Delete result files
	for _, file := range ftJob.ResultFiles {
		tflog.Info(ctx, fmt.Sprintf("Deleting Fine-Tuning Job Result File: %s", file))
		_, err := r.client.Files().DeleteFile(file)
		if err != nil {
			apiError := GetOpenAIAPIError(err)
			if apiError != nil && apiError.HTTPStatusCode == 404 {
				tflog.Info(ctx, "Fine-Tuning Job Result File does not exist")
				err = nil
			}
		}
		if err != nil {
			resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to delete Result File %s, got error: %s", file, err))
			return
		}
	}

	// Delete the fine tuned model
	if ftJob.FineTunedModel != "" {
		tflog.Info(ctx, fmt.Sprintf("Deleting Fine-Tune Model: %s", ftJob.FineTunedModel))
		bDeleted, err := r.client.Models().DeleteFineTuneModel(ftJob.FineTunedModel)
		if err != nil {
			if err, ok := err.(*openai.APIError); ok {
				fmt.Println("openai error:", err.Code)
				// Or whatever other field(s) you need
			}

			resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to delete Fime Tuned Model, got error: %s", err))
			return
		}
		if !bDeleted {
			tflog.Trace(ctx, "Fine Tuned Model not deleted")
		}
		tflog.Trace(ctx, "Fine Tuned Model deleted successfully")
	}
}

func (r *FineTuningJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
