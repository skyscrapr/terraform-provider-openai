package openai

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/skyscrapr/openai-sdk-go/openai"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ModelsDataSource{}

func NewModelsDataSource() datasource.DataSource {
	return &ModelsDataSource{}
}

// ModelsDataSource defines the data source implementation.
type ModelsDataSource struct {
	client *openai.Client
}

// ModelsDataSourceModel describes the data source data model.
type ModelsDataSourceModel struct {
	Id     types.String           `tfsdk:"id"`
	Models []ModelDataSourceModel `tfsdk:"models"`
}

func (d *ModelsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_models"
}

func (d *ModelsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Models data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Models identifier",
				Computed:            true,
			},
			"models": schema.ListNestedAttribute{
				MarkdownDescription: "Models",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: openAIModelAttributes(),
				},
			},
		},
	}
}

func (d *ModelsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openai.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *openai.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *ModelsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ModelsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	models, err := d.client.Models().ListModels()

	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read Models, got error: %s", err))
		return
	}

	for _, v := range models {
		data.Models = append(data.Models, NewModelDataSourceModel(&v))
	}
	data.Id = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
