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
var _ datasource.DataSource = &ModelDataSource{}

func NewFileDataSource() datasource.DataSource {
	return &FileDataSource{}
}

// DataSource defines the data source implementation.
type FileDataSource struct {
	client *openai.Client
}

// ModelDataSourceModel describes the data source data model.
type FileDataSourceModel struct {
	Id       types.String `tfsdk:"id"`
	Bytes    types.Int64  `tfsdk:"bytes"`
	Created  types.Int64  `tfsdk:"created"`
	Filename types.String `tfsdk:"filename"`
	Object   types.String `tfsdk:"object"`
	Purpose  types.String `tfsdk:"fine-tune"`
}

func (d *FileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "file"
}

func (d *FileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "File data source",

		Attributes: openAIFileAttributes(),
	}
}

func (d *FileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FileDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	file, err := d.client.Files().RetrieveFile(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read File, got error: %s", err))
		return
	}

	data = NewFileDataSourceModel(file)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewFileDataSourceModel(f *openai.File) FileDataSourceModel {
	fileDataSourceModel := FileDataSourceModel{
		Id:       types.StringValue(f.Id),
		Bytes:    types.Int64Value(f.Bytes),
		Created:  types.Int64Value(f.CreatedAt),
		Filename: types.StringValue(f.Filename),
		Object:   types.StringValue(f.Object),
		Purpose:  types.StringValue(f.Purpose),
	}
	return fileDataSourceModel
}

func openAIFileAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "File Identifier",
			Required:            true,
		},
		"bytes": schema.Int64Attribute{
			MarkdownDescription: "File size in bytes",
			Computed:            true,
		},
		"created": schema.Int64Attribute{
			MarkdownDescription: "Created Time",
			Computed:            true,
		},
		"filename": schema.StringAttribute{
			MarkdownDescription: "Filename",
			Computed:            true,
		},
		"object": schema.StringAttribute{
			MarkdownDescription: "Object Type",
			Computed:            true,
		},
		"purpose": schema.StringAttribute{
			MarkdownDescription: "Intended use of file. Use 'fine-tune' for Fine-tuning",
			Computed:            true,
		},
	}
}
