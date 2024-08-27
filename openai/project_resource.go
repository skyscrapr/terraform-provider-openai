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
var _ resource.Resource = &ProjectResource{}
var _ resource.ResourceWithImportState = &ProjectResource{}

func NewProjectResource() resource.Resource {
	return &ProjectResource{OpenAIResource: &OpenAIResource{}}
}

// ProjectResource defines the resource implementation.
type ProjectResource struct {
	*OpenAIResource
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Represents an individual project.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The identifier, which can be referenced in API endpoints.",
				Computed:            true,
			},
			"object": schema.StringAttribute{
				MarkdownDescription: "The object type, which is always organization.project",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the project. This appears in reporting.",
				Optional:            true,
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
				MarkdownDescription: "active or archived.",
				Computed:            true,
			},
		},
	}
}

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating Assistant...")

	aReq := openai.ProjectRequest{
		Name: data.Name.ValueStringPointer(),
	}

	project, err := r.client.Projects().CreateProject(&aReq)
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to create project, got error: %s", err))
		return
	}
	tflog.Info(ctx, "Project created successfully")
	data = NewProjectModel(project)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Reading Project with id: %s", data.Id.ValueString()))
	project, err := r.client.Projects().RetrieveProject(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to retrieve project, got error: %s", err))
		return
	}

	data = NewProjectModel(project)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProjectModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	var state ProjectModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating Project: %s", state.Id.ValueString()))

	aReq := openai.ProjectRequest{
		Name: data.Name.ValueStringPointer(),
	}

	project, err := r.client.Projects().ModifyProject(state.Id.ValueString(), aReq)
	if err != nil {
		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to modify project, got error: %s", err))
		return
	}
	tflog.Info(ctx, "Project modified successfully")

	data = NewProjectModel(project)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Archive the project
	tflog.Info(ctx, fmt.Sprintf("Archiving Project: %s", data.Id.ValueString()))
	project, err := r.client.Projects().ArchiveProject(data.Id.ValueString())
	if err != nil {
		if err, ok := err.(*openai.APIError); ok {
			fmt.Println("openai error:", err.Code)
			// Or whatever other field(s) you need
		}

		resp.Diagnostics.AddError("OpenAI Client Error", fmt.Sprintf("Unable to archive project, got error: %s", err))
		return
	}
	if project.Status != "archived" {
		tflog.Trace(ctx, "Project not archived")
	}
	tflog.Trace(ctx, "Project archived successfully")
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
