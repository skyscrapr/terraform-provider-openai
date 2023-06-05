package openai

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/skyscrapr/openai-sdk-go/openai"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ModelDataSource{}

func NewModelDataSource() datasource.DataSource {
	return &ModelDataSource{OpenAIDatasource: &OpenAIDatasource{}}
}

// ModelDataSource defines the data source implementation.
type ModelDataSource struct {
	*OpenAIDatasource
}

// ModelDataSourceModel describes the data source data model.
type ModelDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Created     types.Int64  `tfsdk:"created"`
	Object      types.String `tfsdk:"object"`
	OwnedBy     types.String `tfsdk:"owned_by"`
	Parent      types.String `tfsdk:"parent"`
	Permissions types.List   `tfsdk:"permissions"`
	Root        types.String `tfsdk:"root"`
}

type ModelPermissionModel struct {
	Id                 types.String `tfsdk:"id"`
	Created            types.Int64  `tfsdk:"created"`
	Object             types.String `tfsdk:"object"`
	AllowCreateEngine  types.Bool   `tfsdk:"allow_create_engine"`
	AllowSampling      types.Bool   `tfsdk:"allow_sampling"`
	AllowLogprobs      types.Bool   `tfsdk:"allow_logprobs"`
	AllowSearchIndices types.Bool   `tfsdk:"allow_search_indices"`
	AllowView          types.Bool   `tfsdk:"allow_view"`
	AllowFineTuning    types.Bool   `tfsdk:"allow_fine_tuning"`
	Organization       types.String `tfsdk:"organization"`
	// Group              types.String `tfsdk:"group"`
	IsBlocking types.Bool `tfsdk:"is_blocking"`
}

func (p ModelPermissionModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                   types.StringType,
		"created":              types.Int64Type,
		"object":               types.StringType,
		"allow_create_engine":  types.BoolType,
		"allow_sampling":       types.BoolType,
		"allow_logprobs":       types.BoolType,
		"allow_search_indices": types.BoolType,
		"allow_view":           types.BoolType,
		"allow_fine_tuning":    types.BoolType,
		"organization":         types.StringType,
		// "group":        types.StringType,
		"is_blocking": types.BoolType,
	}
}

func (d *ModelDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_model"
}

func (d *ModelDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Model data source",

		Attributes: openAIModelAttributes(),
	}
}

func (d *ModelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ModelDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	model, err := d.client.Models().RetrieveModel(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read Model, got error: %s", err))
		return
	}

	data = NewModelDataSourceModel(model)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewModelDataSourceModel(model *openai.Model) ModelDataSourceModel {
	modelDataSourceModel := ModelDataSourceModel{
		Id:      types.StringValue(model.ID),
		Created: types.Int64Value(model.CreatedAt),
		Object:  types.StringValue(model.Object),
		OwnedBy: types.StringValue(model.OwnedBy),
		Parent:  types.StringValue(model.Parent),
		Root:    types.StringValue(model.Root),
	}

	// var permissions []types.Object
	var permissions = make([]ModelPermissionModel, len(model.Permission))
	for i, p := range model.Permission {
		permission := ModelPermissionModel{
			Id:                 types.StringValue(p.ID),
			Created:            types.Int64Value(p.CreatedAt),
			Object:             types.StringValue(p.Object),
			AllowCreateEngine:  types.BoolValue(p.AllowCreateEngine),
			AllowSampling:      types.BoolValue(p.AllowSampling),
			AllowLogprobs:      types.BoolValue(p.AllowLogprobs),
			AllowSearchIndices: types.BoolValue(p.AllowSearchIndices),
			AllowView:          types.BoolValue(p.AllowView),
			AllowFineTuning:    types.BoolValue(p.AllowFineTuning),
			Organization:       types.StringValue(p.Organization),
			// Group: types.StringValue(p.Group),
			IsBlocking: types.BoolValue(p.IsBlocking),
		}
		permissions[i] = permission
	}

	modelDataSourceModel.Permissions, _ = types.ListValueFrom(context.TODO(), types.ObjectType{AttrTypes: ModelPermissionModel{}.AttrTypes()}, permissions)
	//modelDataSourceModel.Permissions = permissions
	return modelDataSourceModel
}

func openAIModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "Model Identifier",
			Required:            true,
		},
		"created": schema.Int64Attribute{
			MarkdownDescription: "Created Time",
			Computed:            true,
		},
		"object": schema.StringAttribute{
			MarkdownDescription: "Object Type",
			Computed:            true,
		},
		"owned_by": schema.StringAttribute{
			MarkdownDescription: "Model Owner",
			Computed:            true,
		},
		"parent": schema.StringAttribute{
			MarkdownDescription: "Parent",
			Computed:            true,
		},
		"permissions": schema.ListNestedAttribute{
			MarkdownDescription: "Permissions",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Permission Identifier",
						Computed:            true,
					},
					"created": schema.Int64Attribute{
						MarkdownDescription: "Created Time",
						Computed:            true,
					},
					"object": schema.StringAttribute{
						MarkdownDescription: "Object Type",
						Computed:            true,
					},
					"allow_create_engine": schema.BoolAttribute{
						MarkdownDescription: "Allow Create Engine",
						Computed:            true,
					},
					"allow_sampling": schema.BoolAttribute{
						MarkdownDescription: "Allow Sampling",
						Computed:            true,
					},
					"allow_logprobs": schema.BoolAttribute{
						MarkdownDescription: "Allow Logprobs",
						Computed:            true,
					},
					"allow_search_indices": schema.BoolAttribute{
						MarkdownDescription: "Allow Search Indices",
						Computed:            true,
					},
					"allow_view": schema.BoolAttribute{
						MarkdownDescription: "Allow View",
						Computed:            true,
					},
					"allow_fine_tuning": schema.BoolAttribute{
						MarkdownDescription: "Allow Fine Tuning",
						Computed:            true,
					},
					"organization": schema.StringAttribute{
						MarkdownDescription: "Organization",
						Computed:            true,
					},
					// "group": schema.StringAttribute{
					// 	MarkdownDescription: "Group",
					// 	Computed:            true,
					// },
					"is_blocking": schema.BoolAttribute{
						MarkdownDescription: "Is Blocking",
						Computed:            true,
					},
				},
			},
		},
		"root": schema.StringAttribute{
			MarkdownDescription: "Root",
			Computed:            true,
		},
	}
}
