package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/skyscrapr/openai-sdk-go/openai"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FineTuningJobResource{}
var _ resource.ResourceWithImportState = &FineTuningJobResource{}

func NewAssistantResource() resource.Resource {
	return &AssistantResource{OpenAIResource: &OpenAIResource{}}
}

// AssistantResource defines the resource implementation.
type AssistantResource struct {
	*OpenAIResource
}

func (r *AssistantResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_assistant"
}

func (r *AssistantResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Assistant resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The identifier, which can be referenced in API endpoints.",
				Computed:            true,
			},
			"object": schema.StringAttribute{
				MarkdownDescription: "The object type, which is always assistant.",
				Computed:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) for when the assistant was created.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the assistant. The maximum length is 256 characters.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the assistant. The maximum length is 512 characters.",
				Optional:            true,
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "ID of the model to use. You can use the List models API to see all of your available models.",
				Required:            true,
			},
			"instructions": schema.StringAttribute{
				MarkdownDescription: "The system instructions that the assistant uses. The maximum length is 32768 characters.",
				Optional:            true,
			},
			"tools": schema.ListNestedAttribute{
				MarkdownDescription: "A list of tool enabled on the assistant. There can be a maximum of 128 tools per assistant. Tools can be of types code_interpreter, retrieval, or function.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Tools can be of types code_interpreter, retrieval, or function.",
							Required:            true,
						},
						"function": schema.SingleNestedAttribute{
							MarkdownDescription: "Function definition for tools of type function.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"description": schema.StringAttribute{
									MarkdownDescription: "A description of what the function does, used by the model to choose when and how to call the function.",
									Optional:            true,
								},
								"name": schema.StringAttribute{
									MarkdownDescription: "The name of the function to be called. Must be a-z, A-Z, 0-9, or contain underscores and dashes, with a maximum length of 64.",
									Required:            true,
								},
								"parameters": schema.StringAttribute{
									MarkdownDescription: "The parameters the functions accepts, described as a JSON Schema object.",
									Required:            true,
								},
							},
						},
					},
				},
			},
			"tool_resources": schema.SingleNestedAttribute{
				MarkdownDescription: "A set of resources that are used by the assistant's tools. The resources are specific to the type of tool. For example, the code_interpreter tool requires a list of file IDs, while the file_search tool requires a list of vector store IDs.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"code_interpreter": schema.SingleNestedAttribute{
						MarkdownDescription: "Function definition for tools of type function.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"file_ids": schema.ListAttribute{
								MarkdownDescription: "A list of file IDs attached to this assistant. There can be a maximum of 20 files attached to the assistant. Files are ordered by their creation date in ascending order.",
								ElementType:         types.StringType,
								Optional:            true,
							},
						},
					},
					"file_search": schema.SingleNestedAttribute{
						MarkdownDescription: "Function definition for tools of type function.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"vector_store_ids": schema.ListAttribute{
								MarkdownDescription: "A list of file IDs attached to this assistant. There can be a maximum of 20 files attached to the assistant. Files are ordered by their creation date in ascending order.",
								ElementType:         types.StringType,
								Optional:            true,
							},
							"vector_stores": schema.SingleNestedAttribute{
								MarkdownDescription: "Function definition for tools of type function.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"file_ids": schema.ListAttribute{
										MarkdownDescription: "A list of file IDs attached to this assistant. There can be a maximum of 20 files attached to the assistant. Files are ordered by their creation date in ascending order.",
										ElementType:         types.StringType,
										Optional:            true,
									},
									"metadata": schema.MapAttribute{
										MarkdownDescription: "Set of 16 key-value pairs that can be attached to a vector store. This can be useful for storing additional information about the vector store in a structured format. Keys can be a maximum of 64 characters long and values can be a maxium of 512 characters long.",
										ElementType:         types.StringType,
										Optional:            true,
									},
								},
							},
						},
					},
				},
			},
			"metadata": schema.MapAttribute{
				MarkdownDescription: "Set of 16 key-value pairs that can be attached to an object. This can be useful for storing additional information about the object in a structured format. Keys can be a maximum of 64 characters long and values can be a maxium of 512 characters long.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"temperature": schema.Float64Attribute{
				MarkdownDescription: "What sampling temperature to use, between 0 and 2. Higher values like 0.8 will make the output more random, while lower values like 0.2 will make it more focused and deterministic.",
				Optional:            true,
			},
			"top_p": schema.Float64Attribute{
				MarkdownDescription: "An alternative to sampling with temperature, called nucleus sampling, where the model considers the results of the tokens with top_p probability mass. So 0.1 means only the tokens comprising the top 10% probability mass are considered.",
				Optional:            true,
			},
		},
	}
}

func (r *AssistantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OpenAIAssistantResourceModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating Assistant...")

	aReq := openai.AssistantRequest{
		Model:        data.Model.ValueString(),
		Name:         data.Name.ValueStringPointer(),
		Description:  data.Description.ValueStringPointer(),
		Instructions: data.Instructions.ValueStringPointer(),
		Temperature:  data.Temperature.ValueFloat64(),
		TopP:         data.TopP.ValueFloat64(),
	}

	var toolModels []OpenAIAssistantToolModel
	resp.Diagnostics.Append(data.Tools.ElementsAs(ctx, &toolModels, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var diags diag.Diagnostics
	aReq.Tools, diags = expandAssistantTools(toolModels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	aReq.ToolResources = expandAssistantToolResources(ctx, data.ToolResources)

	assistant, err := r.client.Assistants().CreateAssistant(&aReq)
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to create assistant, got error: %s", err))
		return
	}
	tflog.Info(ctx, "Assistant created successfully")

	data, diags = NewOpenAIAssistantResourceModel(ctx, assistant)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AssistantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OpenAIAssistantResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Reading Assistant with id: %s", data.Id.ValueString()))
	assistant, err := r.client.Assistants().RetrieveAssistant(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to retrieve assistant, got error: %s", err))
		return
	}

	data, diags := NewOpenAIAssistantResourceModel(ctx, assistant)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AssistantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OpenAIAssistantResourceModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Updating Assistant...")

	aReq := openai.AssistantRequest{
		Model:        data.Model.ValueString(),
		Name:         data.Name.ValueStringPointer(),
		Description:  data.Description.ValueStringPointer(),
		Instructions: data.Instructions.ValueStringPointer(),
		Temperature:  data.Temperature.ValueFloat64(),
		TopP:         data.TopP.ValueFloat64(),
	}

	var toolModels []OpenAIAssistantToolModel
	resp.Diagnostics.Append(data.Tools.ElementsAs(ctx, &toolModels, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var diags diag.Diagnostics
	aReq.Tools, diags = expandAssistantTools(toolModels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	aReq.ToolResources = expandAssistantToolResources(ctx, data.ToolResources)

	assistant, err := r.client.Assistants().ModifyAssistant(&aReq)
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to modify assistant, got error: %s", err))
		return
	}
	tflog.Info(ctx, "Assistant modified successfully")

	data, diags = NewOpenAIAssistantResourceModel(ctx, assistant)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AssistantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OpenAIAssistantResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the assistant
	tflog.Info(ctx, fmt.Sprintf("Deleting Assistant: %s", data.Id.ValueString()))
	bDeleted, err := r.client.Assistants().DeleteAssistant(data.Id.ValueString())
	if err != nil {
		if err, ok := err.(*openai.APIError); ok {
			fmt.Println("openai error:", err.Code)
			// Or whatever other field(s) you need
		}

		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to delete assistant, got error: %s", err))
		return
	}
	if !bDeleted {
		tflog.Trace(ctx, "Assistant not deleted")
	}
	tflog.Trace(ctx, "Assistant deleted successfully")
}

func (r *AssistantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func expandAssistantTools(tfList []OpenAIAssistantToolModel) ([]openai.AssistantTool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(tfList) == 0 {
		return nil, diags
	}
	var tools []openai.AssistantTool

	for _, item := range tfList {
		tool := openai.AssistantTool{
			Type: item.Type.ValueString(),
		}
		if item.Function != nil {
			tool.Function = &struct {
				Description *string                "json:\"description,omitempty\""
				Name        string                 "json:\"name\""
				Parameters  map[string]interface{} "json:\"parameters\""
			}{
				Description: item.Function.Description.ValueStringPointer(),
				Name:        *item.Function.Name.ValueStringPointer(),
			}
			if !item.Function.Parameters.IsNull() {
				// Unmarshal the JSON string into the struct
				err := json.Unmarshal([]byte(item.Function.Parameters.ValueString()), &tool.Function.Parameters)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		tools = append(tools, tool)
	}
	return tools, diags
}

func expandAssistantToolResources(ctx context.Context, model *OpenAIAssistantToolResourcesModel) *openai.AssistantToolResources {
	if model == nil {
		return nil
	}
	toolResources := &openai.AssistantToolResources{}
	if !model.CodeInterpreter.IsNull() {
		toolResources.CodeInterpreter = &struct {
			FileIDs []string "json:\"file_ids\""
		}{}
		codeInterpreter := OpenAIAssistantToolResourceCodeInterpreterModel{}
		model.CodeInterpreter.As(ctx, &codeInterpreter, basetypes.ObjectAsOptions{})
		codeInterpreter.FileIDs.ElementsAs(ctx, &toolResources.CodeInterpreter.FileIDs, false)
	}
	if !model.FileSearch.IsNull() {
		toolResources.FileSearch = &struct {
			VectorStoreIDs []string "json:\"vector_store_ids\""
			VectorStores   *struct {
				FileIDs  []string          "json:\"file_ids\""
				MetaData map[string]string "json:\"metadata,omitempty\""
			} "json:\"vector_stores,omitempty\""
		}{}

		fileSearch := OpenAIAssistantToolResourceFileSearchModel{}
		model.FileSearch.As(ctx, &fileSearch, basetypes.ObjectAsOptions{})
		fileSearch.VectorStoreIDs.ElementsAs(ctx, &toolResources.FileSearch.VectorStoreIDs, false)

		if !fileSearch.VectorStores.IsNull() {
			toolResources.FileSearch.VectorStores = &struct {
				FileIDs  []string          "json:\"file_ids\""
				MetaData map[string]string "json:\"metadata,omitempty\""
			}{}
			vectorStores := OpenAIAssistantToolResourceFileSearchVectorStoresModel{}
			fileSearch.VectorStores.As(ctx, &vectorStores, basetypes.ObjectAsOptions{})
			// vectorStore.FileIDs
			// toolResources.FileSearch.VectorStore.FileIDs.ElementsAs(ctx, &toolResources.FileSearch.VectorStore.FileIDs, false)
		}
	}
	// model.CodeInterpreter.As(ctx, toolResources.CodeInterpreter, basetypes.ObjectAsOptions{})
	// model.FileSearch.As(ctx, toolResources.FileSearch, basetypes.ObjectAsOptions{})
	return toolResources
}
