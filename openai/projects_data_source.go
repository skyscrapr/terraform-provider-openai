package openai

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ProjectsDataSource{}

func NewProjectsDataSource() datasource.DataSource {
	return &ProjectsDataSource{OpenAIDatasource: &OpenAIDatasource{}}
}

// ProjectsDataSource defines the data source implementation.
type ProjectsDataSource struct {
	*OpenAIDatasource
}

// ProjectsModel describes the data source data model.
type ProjectsModel struct {
	Id       types.String   `tfsdk:"id"`
	Projects []ProjectModel `tfsdk:"projects"`
}

func (d *ProjectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *ProjectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Projects data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Projects identifier",
				Computed:            true,
			},
			"projects": schema.ListNestedAttribute{
				MarkdownDescription: "Projects",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: openAIProjectAttributes(),
				},
			},
		},
	}
}

func (d *ProjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projects, err := d.client.Projects().ListProjects()

	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read Projects, got error: %s", err))
		return
	}

	for _, v := range projects {
		data.Projects = append(data.Projects, NewProjectModel(&v))
	}
	data.Id = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
