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

func NewFilesDataSource() datasource.DataSource {
	return &FilesDataSource{OpenAIDatasource: &OpenAIDatasource{}}
}

// FilesDataSource defines the data source implementation.
type FilesDataSource struct {
	*OpenAIDatasource
}

// FilesDataSourceModel describes the data source data model.
type FilesDataSourceModel struct {
	Id    types.String      `tfsdk:"id"`
	Files []OpenAIFileModel `tfsdk:"files"`
}

func (d *FilesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_files"
}

func (d *FilesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Files data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Files identifier",
				Computed:            true,
			},
			"files": schema.ListNestedAttribute{
				MarkdownDescription: "Files",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: openAIFileAttributes(),
				},
			},
		},
	}
}

func (d *FilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FilesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	files, err := d.client.Files().ListFiles()

	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read Files, got error: %s", err))
		return
	}

	for _, v := range files {
		data.Files = append(data.Files, NewOpenAIFileModel(&v))
	}
	data.Id = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
