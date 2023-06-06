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
var _ datasource.DataSource = &ModelsDataSource{}

func NewFineTunesDataSource() datasource.DataSource {
	return &FineTunesDataSource{OpenAIDatasource: &OpenAIDatasource{}}
}

// FineTunesDataSource defines the data source implementation.
type FineTunesDataSource struct {
	*OpenAIDatasource
}

// FilesDataSourceModel describes the data source data model.
type FineTunesDataSourceModel struct {
	Id        types.String          `tfsdk:"id"`
	FineTunes []OpenAIFineTuneModel `tfsdk:"finetunes"`
}

func (d *FineTunesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_finetunes"
}

func (d *FineTunesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Files data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Files identifier",
				Computed:            true,
			},
			"finetunes": schema.ListNestedAttribute{
				MarkdownDescription: "Fine Tunes",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: openAIFineTuneAttributes(),
				},
			},
		},
	}
}

func (d *FineTunesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FineTunesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	fineTunes, err := d.client.FineTunes().ListFineTunes()

	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read Fine Tunes, got error: %s", err))
		return
	}

	for _, f := range fineTunes {
		data.FineTunes = append(data.FineTunes, NewOpenAIFineTuneModel(&f))
	}
	data.Id = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
