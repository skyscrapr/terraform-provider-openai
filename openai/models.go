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
	Object   types.String `tfsdk:"object"`
	Purpose  types.String `tfsdk:"purpose"`
}

func (e OpenAIFileModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":       types.StringType,
		"bytes":    types.Int64Type,
		"created":  types.Int64Type,
		"filename": types.StringType,
		"object":   types.StringType,
		"purpose":  types.StringType,
	}
}

func NewOpenAIFileModel(f *openai.File) OpenAIFileModel {
	return OpenAIFileModel{
		Id:       types.StringValue(f.Id),
		Bytes:    types.Int64Value(f.Bytes),
		Created:  types.Int64Value(f.CreatedAt),
		Filename: types.StringValue(f.Filename),
		Object:   types.StringValue(f.Object),
		Purpose:  types.StringValue(f.Purpose),
	}
}

// OpenAIFineTuneModel describes the OpenAI fine-tune model.
type OpenAIFineTuneModel struct {
	Id              types.String `tfsdk:"id"`
	Object          types.String `tfsdk:"object"`
	Model           types.String `tfsdk:"model"`
	Created         types.Int64  `tfsdk:"created"`
	Events          types.List   `tfsdk:"events"`
	FineTunedModel  types.String `tfsdk:"fine_tuned_model"`
	Hyperparams     types.Object `tfsdk:"hyperparams"`
	OrganizationId  types.String `tfsdk:"organization_id"`
	Status          types.String `tfsdk:"status"`
	ResultFiles     types.List   `tfsdk:"result_files"`
	TrainingFiles   types.List   `tfsdk:"training_files"`
	ValidationFiles types.List   `tfsdk:"validation_files"`
	UpdatedAt       types.Int64  `tfsdk:"updated_at"`
}

func (e OpenAIFineTuneModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               types.StringType,
		"object":           types.StringType,
		"model":            types.StringType,
		"created":          types.Int64Type,
		"events":           types.ListType{},
		"fine_tuned_model": types.StringType,
		"hyperparams":      types.ObjectType{},
		"organization_id":  types.StringType,
		"status":           types.StringType,
		"result_files":     types.ListType{},
		"training_files":   types.ListType{},
		"validation_files": types.ListType{},
		"updated_at":       types.Int64Type,
	}
}

func NewOpenAIFineTuneModel(ft *openai.FineTune) OpenAIFineTuneModel {
	ctx := context.TODO()

	fineTuneModel := OpenAIFineTuneModel{
		Id:             types.StringValue(ft.Id),
		Object:         types.StringValue(ft.Object),
		Model:          types.StringValue(ft.Model),
		Created:        types.Int64Value(ft.CreatedAt),
		FineTunedModel: types.StringValue(ft.FineTunedModel),
		OrganizationId: types.StringValue(ft.OrganizationId),
		Status:         types.StringValue(ft.Status),
		UpdatedAt:      types.Int64Value(ft.UpdatedAt),
	}

	hyperparams := OpenAIFineTuneHyperparamsModel{
		BatchSize:              types.Int64Value(ft.Hyperparams.BatchSize),
		LearningRateMultiplier: types.Float64Value(ft.Hyperparams.LearningRateMultiplier),
		NEpochs:                types.Int64Value(ft.Hyperparams.NEpochs),
		PromptLossWeight:       types.Float64Value(ft.Hyperparams.PromptLossWeight),
	}
	fineTuneModel.Hyperparams, _ = types.ObjectValueFrom(ctx, OpenAIFineTuneHyperparamsModel{}.AttrTypes(), hyperparams)

	var events = make([]OpenAIFineTuneEventModel, len(ft.Events))
	for i, e := range ft.Events {
		event := OpenAIFineTuneEventModel{
			Object:  types.StringValue(e.Object),
			Created: types.Int64Value(e.CreatedAt),
			Level:   types.StringValue(e.Level),
			Message: types.StringValue(e.Message),
		}
		events[i] = event
	}
	fineTuneModel.Events, _ = types.ListValueFrom(context.TODO(), types.ObjectType{AttrTypes: OpenAIFineTuneEventModel{}.AttrTypes()}, events)

	var resultFiles = make([]OpenAIFileModel, len(ft.ResultFiles))
	for i, f := range ft.ResultFiles {
		file := NewOpenAIFileModel(&f)
		resultFiles[i] = file
	}
	fineTuneModel.ResultFiles, _ = types.ListValueFrom(context.TODO(), types.ObjectType{AttrTypes: OpenAIFileModel{}.AttrTypes()}, resultFiles)

	var trainingFiles = make([]OpenAIFileModel, len(ft.TrainingFiles))
	for i, f := range ft.TrainingFiles {
		file := NewOpenAIFileModel(&f)
		trainingFiles[i] = file
	}
	fineTuneModel.TrainingFiles, _ = types.ListValueFrom(context.TODO(), types.ObjectType{AttrTypes: OpenAIFileModel{}.AttrTypes()}, trainingFiles)

	var validationFiles = make([]OpenAIFileModel, len(ft.ValidationFiles))
	for i, f := range ft.ValidationFiles {
		file := NewOpenAIFileModel(&f)
		validationFiles[i] = file
	}
	fineTuneModel.ValidationFiles, _ = types.ListValueFrom(context.TODO(), types.ObjectType{AttrTypes: OpenAIFileModel{}.AttrTypes()}, validationFiles)

	return fineTuneModel
}

type OpenAIFineTuneResourceModel struct {
	Id                           types.String `tfsdk:"id"`
	TrainingFile                 types.String `tfsdk:"training_file"`
	ValidationFile               types.String `tfsdk:"validation_file"`
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
	FineTune                     types.Object `tfsdk:"fine_tune"`
}

// func NewOpenAIFineTuneResourceModel(ft *openai.FineTune) OpenAIFineTuneResourceModel {
// 	fineTuneResourceModel := OpenAIFineTuneResourceModel{
// 		FineTune: NewOpenAIFineTuneModel(ft),
// 	}
// 	return fineTuneResourceModel
// }

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

type OpenAIFineTuneHyperparamsModel struct {
	BatchSize              types.Int64   `tfsdk:"batch_size"`
	LearningRateMultiplier types.Float64 `tfsdk:"learning_rate_multiplier"`
	NEpochs                types.Int64   `tfsdk:"n_epochs"`
	PromptLossWeight       types.Float64 `tfsdk:"prompt_loss_weight"`
}

func (e OpenAIFineTuneHyperparamsModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"batch_size":               types.Int64Type,
		"learning_rate_multiplier": types.Float64Type,
		"n_epochs":                 types.Int64Type,
		"prompt_loss_weight":       types.Float64Type,
	}
}
