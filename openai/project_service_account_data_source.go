package openai

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/skyscrapr/openai-sdk-go/openai"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ProjectServiceAccountDataSource{}

func NewProjectServiceAccountDataSource() datasource.DataSource {
	return &ProjectServiceAccountDataSource{OpenAIDatasource: &OpenAIDatasource{}}
}

// ProjectServiceAccountDataSource defines the data source implementation.
type ProjectServiceAccountDataSource struct {
	*OpenAIDatasource
}

// ProjectServiceAccountModel describes the data source data model.
type ProjectServiceAccountModel struct {
	Id         types.String `tfsdk:"id"`
	ProjectId  types.String `tfsdk:"project_id"`
	Object     types.String `tfsdk:"object"`
	Name       types.String `tfsdk:"name"`
	Role       types.String `tfsdk:"role"`
	CreatedAt  types.Int64  `tfsdk:"created_at"`
}

type ProjectServiceAccountResourceModel struct {
	Id         types.String `tfsdk:"id"`
	ProjectId  types.String `tfsdk:"project_id"`
	Object     types.String `tfsdk:"object"`
	Name       types.String `tfsdk:"name"`
	Role       types.String `tfsdk:"role"`
	CreatedAt  types.Int64  `tfsdk:"created_at"`
	ApiKey     *ProjectServiceAccountApiKeyModel `tfsdk:"api_key"`
}

type ProjectServiceAccountApiKeyModel struct {
	Id         types.String `tfsdk:"id"`
	Object     types.String `tfsdk:"object"`
	Value       types.String `tfsdk:"role"`
	CreatedAt  types.Int64  `tfsdk:"created_at"`  
}

func (d *ProjectServiceAccountDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_service_account"
}

func (d *ProjectServiceAccountDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Project Service Account data source",

		Attributes: openAIProjectServiceAccountAttributes(),
	}
}

func (d *ProjectServiceAccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectServiceAccountModel

	// Read Terraform configuration data into the project service account
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectServiceAccount, err := d.client.Projects().RetrieveProjectServiceAccount(data.ProjectId.ValueString(), data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read ProjectServiceAccount, got error: %s", err))
		return
	}

	projectId := data.ProjectId
	data = NewProjectServiceAccountModel(projectServiceAccount)
	data.ProjectId = projectId

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewProjectServiceAccountModel(projectServiceAccount *openai.ProjectServiceAccount) ProjectServiceAccountModel {
	projectServiceAccountModel := ProjectServiceAccountModel{
		Id:         types.StringValue(projectServiceAccount.ID),
		Object:     types.StringValue(projectServiceAccount.Object),
		Name:       types.StringValue(projectServiceAccount.Name),
		Role: 	    types.StringValue(projectServiceAccount.Role),
		CreatedAt:  types.Int64Value(projectServiceAccount.CreatedAt),
	}

	return projectServiceAccountModel
}

func openAIProjectServiceAccountAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The identifier, which can be referenced in API endpoints",
			Required:            true,
		},
		"object": schema.StringAttribute{
			MarkdownDescription: "The object type, which is always organization.project.service_account",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the service account.",
			Computed:            true,
		},
		"role": schema.StringAttribute{
			MarkdownDescription: "owner or member",
			Computed:            true,
		},
		"created_at": schema.Int64Attribute{
			MarkdownDescription: "The Unix timestamp (in seconds) of when the service account was created.",
			Computed:            true,
		},
	}
}
