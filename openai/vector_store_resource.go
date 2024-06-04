package openai

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/skyscrapr/openai-sdk-go/openai"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &VectorStoreResource{}
var _ resource.ResourceWithImportState = &VectorStoreResource{}

func NewVectorStoreResource() resource.Resource {
	return &VectorStoreResource{OpenAIResource: &OpenAIResource{}}
}

// VectorStoreResource defines the resource implementation.
type VectorStoreResource struct {
	*OpenAIResource
}

func (r *VectorStoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vector_store"
}

func (r *VectorStoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "VectorStore resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "VectorStore Identifier",
				Computed:            true,
			},
			"object": schema.StringAttribute{
				MarkdownDescription: "The object type, which is always vector_store.",
				Computed:            true,
			},
			"file_ids": schema.ListAttribute{
				MarkdownDescription: "A list of file IDs attached to this vector store. There can be a maximum of 20 files attached to the assistant. Files are ordered by their creation date in ascending order.",
				ElementType:         types.StringType,
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name",
				Required:            true,
			},
			"expires_after": schema.SingleNestedAttribute{
				MarkdownDescription: "The expiration policy for a vector store.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"anchor": schema.StringAttribute{
						MarkdownDescription: "Anchor timestamp after which the expiration policy applies.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("last_active_at"),
					},
					"days": schema.Int64Attribute{
						MarkdownDescription: "The number of days after the anchor time that the vector store will expire.",
						Required:            true,
					},
				},
			},
			"metadata": schema.MapAttribute{
				MarkdownDescription: "Set of 16 key-value pairs that can be attached to a vector store. This can be useful for storing additional information about the vector store in a structured format. Keys can be a maximum of 64 characters long and values can be a maxium of 512 characters long.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "Created Time",
				Computed:            true,
			},
			"usage_bytes": schema.Int64Attribute{
				MarkdownDescription: "The total number of bytes used by the files in the vector store.",
				Computed:            true,
			},
			"file_counts": schema.SingleNestedAttribute{
				MarkdownDescription: "The expiration policy for a vector store.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"in_progress": schema.Int64Attribute{
						MarkdownDescription: "The number of files that are currently being processed.",
						Computed:            true,
					},
					"completed": schema.Int64Attribute{
						MarkdownDescription: "The number of files that have been successfully processed.",
						Computed:            true,
					},
					"failed": schema.Int64Attribute{
						MarkdownDescription: "The number of files that have failed to process.",
						Computed:            true,
					},
					"cancelled": schema.Int64Attribute{
						MarkdownDescription: "The number of files that were cancelled.",
						Computed:            true,
					},
					"total": schema.Int64Attribute{
						MarkdownDescription: "The total number of files.",
						Computed:            true,
					},
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the vector store, which can be either expired, in_progress, or completed. A status of completed indicates that the vector store is ready for use.",
				Computed:            true,
			},
			"expires_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) for when the vector store will expire.",
				Computed:            true,
			},
			"last_active_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) for when the vector store was last active.",
				Computed:            true,
			},
		},
	}
}

func (r *VectorStoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OpenAIVectorStoreModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vsReq := openai.CreateVectorStoresRequest{
		Name: data.Name.ValueString(),
	}
	resp.Diagnostics.Append(data.FileIDs.ElementsAs(ctx, &vsReq.FileIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(data.Metadata.ElementsAs(ctx, &vsReq.MetaData, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if data.ExpiresAfter != nil {
		vsReq.ExpiresAfter.Anchor = data.ExpiresAfter.Anchor.ValueString()
		vsReq.ExpiresAfter.Days = data.ExpiresAfter.Days.ValueInt64()
	}

	createTimeout := 1 * time.Hour

	vectorStore, err := r.client.VectorStores().CreateVectorStore(&vsReq)
	if err != nil {
		resp.Diagnostics.AddError("CreateVectorStore", fmt.Sprintf("got error: %s", err))
		return
	}

	err = retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		vectorStore, err = r.client.VectorStores().RetrieveVectorStore(vectorStore.Id)
		if err != nil {
			return retry.NonRetryableError(err)
		}
		if vectorStore == nil || vectorStore.FileCounts == nil || vectorStore.FileCounts.InProgress != 0 {
			return retry.RetryableError(fmt.Errorf("file processing still in progress"))
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError("CreateVectorStore", fmt.Sprintf("got error retrieving vector while waiting for files: %s", err))
		return
	}
	if vectorStore.FileCounts.Completed != vectorStore.FileCounts.Total {
		resp.Diagnostics.AddError("CreateVectorStore", "Failed to process all files")
		return
	}
	tflog.Info(ctx, "Vector Store created successfully")

	var diags diag.Diagnostics
	data, diags = NewOpenAIVectoreStoreModel(ctx, vectorStore, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VectorStoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OpenAIVectorStoreModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vectorStore, err := r.client.VectorStores().RetrieveVectorStore(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("RetrieveVectorStore", fmt.Sprintf("got error: %s", err))
		return
	}

	var diags diag.Diagnostics
	data, diags = NewOpenAIVectoreStoreModel(ctx, vectorStore, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VectorStoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "Update not supported.")
}

func (r *VectorStoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OpenAIVectorStoreModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	deletionStatus, err := r.client.VectorStores().DeleteVectorStore(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("DeleteVectorStore", fmt.Sprintf("got error: %s", err))
		return
	}

	if !deletionStatus.Deleted {
		resp.Diagnostics.AddError("DeleteVectorStore", "unknown error. Could not delete vector store")
		return
	}
}

func (r *VectorStoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
