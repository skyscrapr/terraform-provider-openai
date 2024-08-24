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
var _ datasource.DataSource = &ProjectDataSource{}

func NewProjectDataSource() datasource.DataSource {
	return &ProjectDataSource{OpenAIDatasource: &OpenAIDatasource{}}
}

// ProjectDataSource defines the data source implementation.
type ProjectDataSource struct {
	*OpenAIDatasource
}

// ProjectModel describes the data source data model.
type ProjectModel struct {
	Id         types.String `tfsdk:"id"`
	Object     types.String `tfsdk:"object"`
	Name       types.String `tfsdk:"name"`
	CreatedAt  types.Int64  `tfsdk:"created_at"`
	ArchivedAt types.Int64  `tfsdk:"archived_at"`
	Status     types.String `tfsdk:"status"`
}

func (d *ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *ProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Project data source",

		Attributes: openAIProjectAttributes(),
	}
}

func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectModel

	// Read Terraform configuration data into the project
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	project, err := d.client.Projects().RetrieveProject(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read Project, got error: %s", err))
		return
	}

	data = NewProjectModel(project)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewProjectModel(project *openai.Project) ProjectModel {
	projectModel := ProjectModel{
		Id:         types.StringValue(project.ID),
		Object:     types.StringValue(project.Object),
		Name:       types.StringValue(project.Name),
		CreatedAt:  types.Int64Value(project.CreatedAt),
		ArchivedAt: types.Int64Value(project.ArchivedAt),
		Status:     types.StringValue(project.Status),
	}

	return projectModel
}

func openAIProjectAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The identifier, which can be referenced in API endpoints",
			Required:            true,
		},
		"object": schema.StringAttribute{
			MarkdownDescription: "The object type, which is always organization.project",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the project. This appears in reporting.",
			Computed:            true,
		},
		"created_at": schema.Int64Attribute{
			MarkdownDescription: "The Unix timestamp (in seconds) of when the project was created.",
			Computed:            true,
		},
		"archived_at": schema.Int64Attribute{
			MarkdownDescription: "The Unix timestamp (in seconds) of when the project was archived or null.",
			Computed:            true,
		},
		"status": schema.StringAttribute{
			MarkdownDescription: "active or archived",
			Computed:            true,
		},
	}
}
