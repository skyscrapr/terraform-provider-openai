package openai

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/skyscrapr/openai-sdk-go/openai"
)

// OpenAIFileModel describes the OpenAI file model.
type OpenAIFileModel struct {
	Id       types.String `tfsdk:"id"`
	Bytes    types.Int64  `tfsdk:"bytes"`
	Created  types.Int64  `tfsdk:"created"`
	Filename types.String `tfsdk:"filename"`
	Filepath types.String `tfsdk:"filepath"`
	Object   types.String `tfsdk:"object"`
	Purpose  types.String `tfsdk:"purpose"`
}

func (e OpenAIFileModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":       types.StringType,
		"bytes":    types.Int64Type,
		"created":  types.Int64Type,
		"filename": types.StringType,
		"filepath": types.StringType,
		"object":   types.StringType,
		"purpose":  types.StringType,
	}
}

func NewOpenAIFileModelWithPath(f *openai.File, path string) OpenAIFileModel {
	return OpenAIFileModel{
		Id:       types.StringValue(f.Id),
		Bytes:    types.Int64Value(f.Bytes),
		Created:  types.Int64Value(f.CreatedAt),
		Filename: types.StringValue(f.Filename),
		Filepath: types.StringValue(path),
		Object:   types.StringValue(f.Object),
		Purpose:  types.StringValue(f.Purpose),
	}
}

func NewOpenAIFileModel(f *openai.File) OpenAIFileModel {
	return NewOpenAIFileModelWithPath(f, f.Filename)
}

// OpenAIFineTuningJobModel describes the OpenAI fine-tuning job model.
type OpenAIFineTuningJobModel struct {
	Id             types.String `tfsdk:"id"`
	Object         types.String `tfsdk:"object"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	FinishedAt     types.Int64  `tfsdk:"finished_at"`
	Model          types.String `tfsdk:"model"`
	FineTunedModel types.String `tfsdk:"fine_tuned_model"`
	OrganizationId types.String `tfsdk:"organization_id"`
	Status         types.String `tfsdk:"status"`
	Hyperparams    types.Object `tfsdk:"hyperparams"`
	TrainingFile   types.String `tfsdk:"training_file"`
	ValidationFile types.String `tfsdk:"validation_file"`
	ResultFiles    types.List   `tfsdk:"result_files"`
	TrainedTokens  types.Int64  `tfsdk:"trained_tokens"`
	Suffix         types.String `tfsdk:"suffix"`
	Wait           types.Bool   `tfsdk:"wait"`
}

func (e OpenAIFineTuningJobModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               types.StringType,
		"object":           types.StringType,
		"created_at":       types.Int64Type,
		"finished_at":      types.Int64Type,
		"model":            types.StringType,
		"fine_tuned_model": types.StringType,
		"organization_id":  types.StringType,
		"status":           types.StringType,
		"hyperparams":      types.ObjectType{AttrTypes: OpenAIFineTuningJobHyperparamsModel{}.AttrTypes()},
		"training_file":    types.StringType,
		"validation_file":  types.StringType,
		"result_files":     types.ListType{ElemType: types.ObjectType{AttrTypes: OpenAIFileModel{}.AttrTypes()}},
		"trained_tokens":   types.Int64Type,
	}
}

func NewOpenAIFineTuningJobModel(ft *openai.FineTuningJob) OpenAIFineTuningJobModel {
	ctx := context.TODO()

	ftJobModel := OpenAIFineTuningJobModel{
		Id:             types.StringValue(ft.Id),
		Object:         types.StringValue(ft.Object),
		CreatedAt:      types.Int64Value(ft.CreatedAt),
		FinishedAt:     types.Int64Value(ft.FinishedAt),
		Model:          types.StringValue(ft.Model),
		FineTunedModel: types.StringValue(ft.FineTunedModel),
		OrganizationId: types.StringValue(ft.OrganizationId),
		Status:         types.StringValue(ft.Status),
		TrainingFile:   types.StringValue(ft.TrainingFile),
		TrainedTokens:  types.Int64Value(ft.TrainedTokens),
	}

	if ft.ValidationFile != nil {
		ftJobModel.ValidationFile = types.StringValue(*ft.ValidationFile)
	}

	ftJobModel.Hyperparams, _ = types.ObjectValueFrom(ctx, OpenAIFineTuningJobHyperparamsModel{}.AttrTypes(), ft.Hyperparams)
	ftJobModel.ResultFiles, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: OpenAIFileModel{}.AttrTypes()}, ft.ResultFiles)

	return ftJobModel
}

type OpenAIFineTuneEventModel struct {
	Object  types.String `tfsdk:"object"`
	Created types.Int64  `tfsdk:"created"`
	Level   types.String `tfsdk:"level"`
	Message types.String `tfsdk:"message"`
}

func (e OpenAIFineTuneEventModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"created": types.Int64Type,
		"object":  types.StringType,
		"level":   types.StringType,
		"message": types.StringType,
	}
}

type OpenAIFineTuningJobHyperparamsModel struct {
	NEpochs types.Int64 `tfsdk:"n_epochs"`
}

func (e OpenAIFineTuningJobHyperparamsModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"n_epochs": types.StringType,
	}
}
