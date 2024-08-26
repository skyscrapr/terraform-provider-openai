package openai

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/skyscrapr/openai-sdk-go/openai"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ProjectServiceAccountResource{}
var _ resource.ResourceWithImportState = &ProjectServiceAccountResource{}

func NewProjectServiceAccountResource() resource.Resource {
	return &ProjectServiceAccountResource{OpenAIResource: &OpenAIResource{}}
}

// ProjectServiceAccountResource defines the resource implementation.
type ProjectServiceAccountResource struct {
	*OpenAIResource
}

func (r *ProjectServiceAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_service_account"
}

func (r *ProjectServiceAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Represents an individual project service account.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The identifier, which can be referenced in API endpoints.",
				Computed:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The identifier, which can be referenced in API endpoints.",
				Required:            true,
			},
			"object": schema.StringAttribute{
				MarkdownDescription: "The object type, which is always organization.project.service_account",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the service account.",
				Optional:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "owner or member",
				Optional:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) of when the project was created.",
				Computed:            true,
			},
			"api_key": schema.SingleNestedAttribute{
				MarkdownDescription: "A list of tool enabled on the assistant. There can be a maximum of 128 tools per assistant. Tools can be of types code_interpreter, retrieval, or function.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "The identifier, which can be referenced in API endpoints.",
						Computed:            true,
					},
					"object": schema.StringAttribute{
						MarkdownDescription: "The object type, which is always organization.project.service_account.api_key.",
						Computed:            true,
					},
					"value": schema.StringAttribute{
						MarkdownDescription: "The value of the api_key secret.",
						Computed:            true,
						Sensitive: true,
					},
					"created_at": schema.Int64Attribute{
						MarkdownDescription: "The Unix timestamp (in seconds) of when the api_key was created.",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (r *ProjectServiceAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectServiceAccountResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating Assistant...")

	aReq := openai.ProjectServiceAccountRequest{
		ProjectID: data.ProjectId.ValueString(),
		Name: data.Name.ValueStringPointer(),
	}

	projectServiceAccount, err := r.client.Projects().CreateProjectServiceAccount(&aReq)
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to create project service account, got error: %s", err))
		return
	}
	tflog.Info(ctx, "Project Service Account created successfully")
	data = NewProjectServiceAccountModel(projectServiceAccount)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectServiceAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectServiceAccountModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Reading ProjectServiceAccount with id: %s", data.Id.ValueString()))
	projectServiceAccount, err := r.client.Projects().RetrieveProjectServiceAccount(data.ProjectId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to retrieve project service account, got error: %s", err))
		return
	}

	data = NewProjectServiceAccountModel(projectServiceAccount)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectServiceAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Project Service Account is not supported")
}

func (r *ProjectServiceAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectServiceAccountModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the project service account
	tflog.Info(ctx, fmt.Sprintf("Deleting Project Service Account: %s", data.Id.ValueString()))
	bDeleted, err := r.client.Projects().DeleteProjectServiceAccount(data.ProjectId.ValueString(), data.Id.ValueString())
	if err != nil {
		if err, ok := err.(*openai.APIError); ok {
			fmt.Println("openai error:", err.Code)
			// Or whatever other field(s) you need
		}

		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to delete project service account, got error: %s", err))
		return
	}
	if !bDeleted {
		tflog.Trace(ctx, "Project Service Account not deleted")
	}
	tflog.Trace(ctx, "Project Service Account deleted successfully")
}

func (r *ProjectServiceAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
