package openai

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/skyscrapr/openai-sdk-go/openai"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FileResource{}
var _ resource.ResourceWithImportState = &FileResource{}

func NewFileResource() resource.Resource {
	return &FileResource{OpenAIResource: &OpenAIResource{}}
}

// FileResource defines the resource implementation.
type FileResource struct {
	*OpenAIResource
}

func (r *FileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (r *FileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "File resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "File Identifier",
				Computed:            true,
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
			"filepath": schema.StringAttribute{
				MarkdownDescription: "Filename",
				Required:            true,
			},
			"object": schema.StringAttribute{
				MarkdownDescription: "Object Type",
				Computed:            true,
			},
			"purpose": schema.StringAttribute{
				MarkdownDescription: "Intended use of file. Use 'fine-tune' for Fine-tuning",
				Computed:            true,
				Default:             stringdefault.StaticString("fine-tune"),
			},
			// "timeouts": Timeouts(ctx, timeouts.Opts{
			//     Delete: true,
			// }),
		},
		// Blocks: map[string]schema.Block{
		// 	"timeouts": Timeouts(ctx, timeouts.Opts{
		//         Delete: true,
		//     }),
		// },
	}
}

func (r *FileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OpenAIFileModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	filePath, err := GetFilePath(data.Filepath.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("UploadFile", fmt.Sprintf("Unable to upload %s, got error: %s", *filePath, err))
		return
	}

	file, err := r.client.Files().UploadFile(&openai.UploadFileRequest{
		File:    *filePath,
		Purpose: data.Purpose.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to upload File: %s", err))
		return
	}
	tflog.Trace(ctx, "Uploaded file successfully")

	data = NewOpenAIFileModelWithPath(file, data.Filepath.ValueString())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OpenAIFileModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	file, err := r.client.Files().RetrieveFile(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to read File, got error: %s", err))
		return
	}

	data = NewOpenAIFileModelWithPath(file, data.Filepath.ValueString())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "Update not supported.")
}

func (r *FileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OpenAIFileModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// HACK: Could not implement timeouts due to error
	destroyTimeout := 20 * time.Minute

	err := retry.RetryContext(ctx, destroyTimeout, func() *retry.RetryError {
		bDeleted, err := r.client.Files().DeleteFile(data.Id.ValueString())
		if err != nil {
			if openaierr, ok := err.(*openai.APIError); ok {
				if openaierr.HTTPStatusCode == 409 {
					tflog.Info(ctx, fmt.Sprintf("%s - Retrying...", err))
					return retry.RetryableError(err)
				}
			}
			return retry.NonRetryableError(err)
		}
		if bDeleted {
			tflog.Trace(ctx, "File deleted successfully")
		} else {
			tflog.Trace(ctx, "File not deleted")
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to delete File, got error: %s", err))
		return
	}
}

func (r *FileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func openAIFileResourceAttributes() map[string]schema.Attribute {
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
		"filepath": schema.StringAttribute{
			MarkdownDescription: "Filepath",
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
