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
var _ datasource.DataSource = &ProjectServiceAccountsDataSource{}

func NewProjectServiceAccountsDataSource() datasource.DataSource {
	return &ProjectServiceAccountsDataSource{OpenAIDatasource: &OpenAIDatasource{}}
}

// ProjectServiceAccountsDataSource defines the data source implementation.
type ProjectServiceAccountsDataSource struct {
	*OpenAIDatasource
}

// ProjectServiceAccountsModel describes the data source data model.
type ProjectServiceAccountsModel struct {
	Id       types.String   `tfsdk:"id"`
	ProjectServiceAccounts []ProjectServiceAccountModel `tfsdk:"project_service_accounts"`
}

func (d *ProjectServiceAccountsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_service_accounts"
}

func (d *ProjectServiceAccountsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Project Service Accounts data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Projects identifier",
				Computed:            true,
			},
			"project_service_accounts": schema.ListNestedAttribute{
				MarkdownDescription: "Project Service Accounts",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: openAIProjectServiceAccountAttributes(),
				},
			},
		},
	}
}

func (d *ProjectServiceAccountsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectServiceAccountsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectServiceAccounts, err := d.client.Projects().ListProjectServiceAccounts()

	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read Project Service Accounts, got error: %s", err))
		return
	}

	for _, v := range projectServiceAccounts {
		data.ProjectServiceAccounts = append(data.ProjectServiceAccounts, NewProjectServiceAccountModel(&v))
	}
	data.Id = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
