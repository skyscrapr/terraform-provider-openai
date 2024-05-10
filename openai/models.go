package openai

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
		"result_files":     types.ListType{ElemType: types.StringType},
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

	h := OpenAIFineTuningJobHyperparamsModel{
		NEpochs: types.Int64Value(ft.Hyperparams.NEpochs),
	}
	ftJobModel.Hyperparams, _ = types.ObjectValueFrom(ctx, OpenAIFineTuningJobHyperparamsModel{}.AttrTypes(), h)

	ftJobModel.ResultFiles, _ = types.ListValueFrom(ctx, types.StringType, ft.ResultFiles)

	return ftJobModel
}

type OpenAIFineTuningJobResourceModel struct {
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

func NewOpenAIFineTuningJobResourceModel(ft *openai.FineTuningJob, wait bool) OpenAIFineTuningJobResourceModel {
	ctx := context.TODO()

	ftJobModel := OpenAIFineTuningJobResourceModel{
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
		Suffix:         types.StringValue(""),
		Wait:           types.BoolValue(wait),
	}

	if ft.ValidationFile != nil {
		ftJobModel.ValidationFile = types.StringValue(*ft.ValidationFile)
	}

	h := OpenAIFineTuningJobHyperparamsModel{
		NEpochs: types.Int64Value(ft.Hyperparams.NEpochs),
	}
	ftJobModel.Hyperparams, _ = types.ObjectValueFrom(ctx, OpenAIFineTuningJobHyperparamsModel{}.AttrTypes(), h)

	ftJobModel.ResultFiles, _ = types.ListValueFrom(ctx, types.StringType, ft.ResultFiles)

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
		"n_epochs": types.Int64Type,
	}
}

type OpenAIAssistantResourceModel struct {
	Id            types.String                       `tfsdk:"id"`
	Object        types.String                       `tfsdk:"object"`
	CreatedAt     types.Int64                        `tfsdk:"created_at"`
	Name          types.String                       `tfsdk:"name"`
	Description   types.String                       `tfsdk:"description"`
	Model         types.String                       `tfsdk:"model"`
	Instructions  types.String                       `tfsdk:"instructions"`
	Tools         types.List                         `tfsdk:"tools"`
	ToolResources *OpenAIAssistantToolResourcesModel `tfsdk:"tool_resources"`
	Metadata      types.Map                          `tfsdk:"metadata"`
	Temperature   types.Float64                      `tfsdk:"temperature"`
	TopP          types.Float64                      `tfsdk:"top_p"`
}

func (e OpenAIAssistantResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":             types.StringType,
		"object":         types.StringType,
		"created_at":     types.Int64Type,
		"name":           types.StringType,
		"description":    types.StringType,
		"model":          types.StringType,
		"instructions":   types.StringType,
		"tools":          types.ListType{ElemType: types.ObjectType{AttrTypes: OpenAIAssistantToolModel{}.AttrTypes()}},
		"tool_resources": types.ObjectType{AttrTypes: OpenAIAssistantToolResourcesModel{}.AttrTypes()},
		"metadata":       types.MapType{ElemType: types.StringType},
		"temperature":    types.Float64Type,
		"top_p":          types.Float64Type,
	}
}

func NewOpenAIAssistantResourceModel(ctx context.Context, assistant *openai.Assistant) (OpenAIAssistantResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := OpenAIAssistantResourceModel{}
	model.Id = types.StringValue(assistant.Id)
	model.Object = types.StringValue(assistant.Object)
	model.Name = types.StringPointerValue(assistant.Name)
	model.Description = types.StringPointerValue(assistant.Description)
	model.Model = types.StringValue(assistant.Model)
	model.Instructions = types.StringPointerValue(assistant.Instructions)

	if len(assistant.MetaData) == 0 {
		model.Metadata = types.MapNull(types.StringType)
	} else {
		model.Metadata, _ = types.MapValueFrom(ctx, types.StringType, assistant.MetaData)
	}

	var tools = make([]OpenAIAssistantToolModel, len(assistant.Tools))
	for i, t := range assistant.Tools {
		tool := OpenAIAssistantToolModel{
			Type: types.StringValue(t.Type),
		}

		if t.Function == nil {
			tool.Function = types.ObjectNull(OpenAIAssistantToolFunctionModel{}.AttrTypes())
		} else {
			f := OpenAIAssistantToolFunctionModel{
				Name:       types.StringValue(t.Function.Name),
				Parameters: types.StringValue(t.Function.Parameters),
			}
			if t.Function.Description != nil {
				f.Description = types.StringPointerValue(t.Function.Description)
			}
			tool.Function, _ = types.ObjectValueFrom(ctx, OpenAIAssistantToolFunctionModel{}.AttrTypes(), f)
		}
		tools[i] = tool
	}
	model.Tools, diags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: OpenAIAssistantToolModel{}.AttrTypes()}, tools)
	if diags.HasError() {
		return model, diags
	}

	if assistant.ToolResources != nil {
		model.ToolResources = &OpenAIAssistantToolResourcesModel{}
		if assistant.ToolResources.CodeInterpreter == nil {
			model.ToolResources.CodeInterpreter = types.ObjectNull(OpenAIAssistantToolResourceCodeInterpreterModel{}.AttrTypes())
		} else {
			codeInterpreter := &OpenAIAssistantToolResourceCodeInterpreterModel{}
			codeInterpreter.FileIDs, diags = types.ListValueFrom(ctx, types.StringType, assistant.ToolResources.CodeInterpreter.FileIDs)
			if diags.HasError() {
				return model, diags
			}
			model.ToolResources.CodeInterpreter, diags = types.ObjectValueFrom(ctx, OpenAIAssistantToolResourceCodeInterpreterModel{}.AttrTypes(), codeInterpreter)
		}
		if assistant.ToolResources.FileSearch == nil {
			model.ToolResources.FileSearch = types.ObjectNull(OpenAIAssistantToolResourceFileSearchModel{}.AttrTypes())
		} else {
			fileSearch := &OpenAIAssistantToolResourceFileSearchModel{}
			fileSearch.VectorStoreIDs, diags = types.ListValueFrom(ctx, types.StringType, assistant.ToolResources.FileSearch.VectorStoreIDs)
			if diags.HasError() {
				return model, diags
			}
			if assistant.ToolResources.FileSearch.VectorStores == nil {
				fileSearch.VectorStores = types.ObjectNull(OpenAIAssistantToolResourceFileSearchVectorStoresModel{}.AttrTypes())
			} else {
				vectorStores := &OpenAIAssistantToolResourceFileSearchVectorStoresModel{}
				vectorStores.FileIDs, diags = types.ListValueFrom(ctx, types.StringType, assistant.ToolResources.FileSearch.VectorStores.FileIDs)
				if diags.HasError() {
					return model, diags
				}
				vectorStores.MetaData, diags = types.MapValueFrom(ctx, types.StringType, assistant.ToolResources.FileSearch.VectorStores.MetaData)
				if diags.HasError() {
					return model, diags
				}
			}
			model.ToolResources.FileSearch, diags = types.ObjectValueFrom(ctx, OpenAIAssistantToolResourceFileSearchModel{}.AttrTypes(), fileSearch)
		}
	}
	if diags.HasError() {
		return model, diags
	}

	return model, diags
}

type OpenAIAssistantToolModel struct {
	Type     types.String `tfsdk:"type"`
	Function types.Object `tfsdk:"function"`
}

func (e OpenAIAssistantToolModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":     types.StringType,
		"function": types.ObjectType{AttrTypes: OpenAIAssistantToolFunctionModel{}.AttrTypes()},
	}
}

type OpenAIAssistantToolResourcesModel struct {
	CodeInterpreter types.Object `tfsdk:"code_interpreter"`
	FileSearch      types.Object `tfsdk:"file_search"`
}

func (e OpenAIAssistantToolResourcesModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"code_interpreter": types.ObjectType{AttrTypes: OpenAIAssistantToolResourceCodeInterpreterModel{}.AttrTypes()},
		"file_search":      types.ObjectType{AttrTypes: OpenAIAssistantToolResourceFileSearchModel{}.AttrTypes()},
	}
}

type OpenAIAssistantToolResourceCodeInterpreterModel struct {
	FileIDs types.List `tfsdk:"file_ids"`
}

func (e OpenAIAssistantToolResourceCodeInterpreterModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"file_ids": types.ListType{ElemType: types.StringType},
	}
}

type OpenAIAssistantToolResourceFileSearchModel struct {
	VectorStoreIDs types.List   `tfsdk:"vector_store_ids"`
	VectorStores   types.Object `tfsdk:"vector_stores"`
}

func (e OpenAIAssistantToolResourceFileSearchModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"vector_store_ids": types.ListType{ElemType: types.StringType},
		"vector_stores":    types.ObjectType{AttrTypes: OpenAIAssistantToolResourceFileSearchVectorStoresModel{}.AttrTypes()},
	}
}

type OpenAIAssistantToolResourceFileSearchVectorStoresModel struct {
	FileIDs  types.List `tfsdk:"file_ids"`
	MetaData types.Map  `tfsdk:"metadata"`
}

func (e OpenAIAssistantToolResourceFileSearchVectorStoresModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"file_ids": types.ListType{ElemType: types.StringType},
		"metadata": types.MapType{ElemType: types.StringType},
	}
}

type OpenAIAssistantToolFunctionModel struct {
	Description types.String `tfsdk:"description"`
	Name        types.String `tfsdk:"name"`
	Parameters  types.String `tfsdk:"parameters"`
}

func (e OpenAIAssistantToolFunctionModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"description": types.StringType,
		"name":        types.StringType,
		"parameters":  types.StringType,
	}
}
