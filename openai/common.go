package openai

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/skyscrapr/openai-sdk-go/openai"
)

type OpenAIDatasource struct {
	client *openai.Client
}

func (d *OpenAIDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openai.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Client Type",
			fmt.Sprintf("Expected *openai.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

type OpenAIResource struct {
	client *openai.Client
}

func (d *OpenAIResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openai.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Client Type",
			fmt.Sprintf("Expected *openai.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// OpenAIFileModel describes the OpenAI file model.
type OpenAIFileModel struct {
	Id       types.String `tfsdk:"id"`
	Bytes    types.Int64  `tfsdk:"bytes"`
	Created  types.Int64  `tfsdk:"created"`
	Filename types.String `tfsdk:"filename"`
	Object   types.String `tfsdk:"object"`
	Purpose  types.String `tfsdk:"purpose"`
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
	Id              types.String      `tfsdk:"id"`
	Object          types.String      `tfsdk:"object"`
	Model           types.String      `tfsdk:"model"`
	Created         types.Int64       `tfsdk:"created"`
	Events          types.List        `tfsdk:"events"`
	FineTunedModel  types.String      `tfsdk:"fine_tuned_model"`
	Hyperparams     types.List        `tfsdk:"hyperparams"`
	OrganizationId  types.String      `tfsdk:"organization_id"`
	Status          types.String      `tfsdk:"status"`
	ResultFiles     []OpenAIFileModel `tfsdk:"result_files"`
	TrainingFiles   []OpenAIFileModel `tfsdk:"training_files"`
	ValidationFiles []OpenAIFileModel `tfsdk:"validation_files"`
	UpdatedAt       types.Int64       `tfsdk:"updated_at"`
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

type OpenAIFineTuneHyperparamModel struct {
	BatchSize              types.Int64   `tfsdk:"batch_size"`
	LearningRateMultiplier types.Float64 `tfsdk:"learning_rate_multiplier"`
	NEpochs                types.Int64   `tfsdk:"n_epochs"`
	PromptLossWeight       types.Float64 `tfsdk:"prompt_loss_weight"`
}

func (e OpenAIFineTuneHyperparamModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"batch_size":               types.Int64Type,
		"learning_rate_multiplier": types.Float64Type,
		"n_epochs":                 types.Int64Type,
		"prompt_loss_weight":       types.Float64Type,
	}
}

func NewOpenAIFineTuneModel(ft *openai.FineTune) OpenAIFineTuneModel {
	fineTuneDatasourceModel := OpenAIFineTuneModel{
		Id:             types.StringValue(ft.Id),
		Object:         types.StringValue(ft.Object),
		Model:          types.StringValue(ft.Model),
		Created:        types.Int64Value(ft.CreatedAt),
		FineTunedModel: types.StringValue(ft.FineTunedModel),
		OrganizationId: types.StringValue(ft.OrganizationId),
		Status:         types.StringValue(ft.Status),
		UpdatedAt:      types.Int64Value(ft.CreatedAt),
	}

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
	fineTuneDatasourceModel.Events, _ = types.ListValueFrom(context.TODO(), types.ObjectType{AttrTypes: OpenAIFineTuneEventModel{}.AttrTypes()}, events)

	var hyperparams = make([]OpenAIFineTuneHyperparamModel, len(ft.Hyperparams))
	for i, h := range ft.Hyperparams {
		hyperparam := OpenAIFineTuneHyperparamModel{
			BatchSize:              types.Int64Value(h.BatchSize),
			LearningRateMultiplier: types.Float64Value(float64(h.LearningRateMultiplier)),
			NEpochs:                types.Int64Value(h.NEpochs),
			PromptLossWeight:       types.Float64Value(float64(h.PromptLossWeight)),
		}
		hyperparams[i] = hyperparam
	}
	fineTuneDatasourceModel.Hyperparams, _ = types.ListValueFrom(context.TODO(), types.ObjectType{AttrTypes: OpenAIFineTuneHyperparamModel{}.AttrTypes()}, hyperparams)

	var resultFiles = make([]OpenAIFileModel, len(ft.ResultFiles))
	for i, f := range ft.ResultFiles {
		file := NewOpenAIFileModel(&f)
		resultFiles[i] = file
	}
	fineTuneDatasourceModel.ResultFiles = resultFiles

	var trainingFiles = make([]OpenAIFileModel, len(ft.TrainingFiles))
	for i, f := range ft.TrainingFiles {
		file := NewOpenAIFileModel(&f)
		trainingFiles[i] = file
	}
	fineTuneDatasourceModel.TrainingFiles = trainingFiles

	var validationFiles = make([]OpenAIFileModel, len(ft.ValidationFiles))
	for i, f := range ft.ValidationFiles {
		file := NewOpenAIFileModel(&f)
		validationFiles[i] = file
	}
	fineTuneDatasourceModel.ValidationFiles = validationFiles

	return fineTuneDatasourceModel
}

func GetFilePath(filePath string) (*string, error) {
	// if filepath.IsAbs(filePath) {
	// 	return &filePath, nil
	// }

	// tfPath := os.Getenv("TF_DATA_DIR")
	// if tfPath == "" {
	// 	return nil, fmt.Errorf("TF_DATA_DIR environment variable not set")
	// }
	// newFilePath := filepath.Join(tfPath, filePath)
	// // _, err := os.Stat(newFilePath)
	// return &newFilePath, nil

	configDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	absFilePath := filepath.Join(configDir, filePath)
	return &absFilePath, nil
}
