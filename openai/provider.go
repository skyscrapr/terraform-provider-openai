package openai

import (
	"context"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/skyscrapr/openai-sdk-go/openai"
)

// Ensure OpenAIProvider satisfies various provider interfaces.
var _ provider.Provider = &OpenAIProvider{}

// OpenAIProvider defines the provider implementation.
type OpenAIProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// OpenAIProviderModel describes the provider data model.
type OpenAIProviderModel struct {
	ApiKey         types.String `tfsdk:"api_key"`
	AdminKey       types.String `tfsdk:"admin_key"`
	BaseURL        types.String `tfsdk:"base_url"`
	OrganizationID types.String `tfsdk:"organization_id"`
}

func (p *OpenAIProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "openai"
	resp.Version = p.version
}

func (p *OpenAIProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"admin_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"base_url": schema.StringAttribute{
				Optional: true,
			},
			"organization_id": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *OpenAIProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data OpenAIProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if data.ApiKey.IsUnknown() && data.AdminKey.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown OpenAI Authorization Tokens",
			"The provider cannot create the OpenAI API client as there is an unknown configuration value for the OpenAI API authorization token. "+
				"Either set the value statically in the configuration, or use the OPENAI_API_KEY or OPENAI_ADMIN_KEY environment variables.",
		)
		return
	}

	client := configureClient(data)

	// Make the OpenAI client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *OpenAIProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAssistantResource,
		NewFileResource,
		NewFineTuningJobResource,
		NewProjectResource,
		NewVectorStoreResource,
		NewProjectServiceAccountResource,
	}
}

func (p *OpenAIProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFilesDataSource,
		NewFileDataSource,
		NewFineTuningJobsDataSource,
		NewFineTuningJobDataSource,
		NewModelsDataSource,
		NewModelDataSource,
		NewProjectsDataSource,
		NewProjectDataSource,
		NewProjectServiceAccountsDataSource,
		NewProjectServiceAccountDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpenAIProvider{
			version: version,
		}
	}
}

func configureClient(data OpenAIProviderModel) (*openai.Client) {
	api_key := os.Getenv("OPENAI_API_KEY")
	if !data.ApiKey.IsNull() {
		api_key = data.ApiKey.ValueString()
	}

	admin_key := os.Getenv("OPENAI_ADMIN_KEY")
	if !data.AdminKey.IsNull() {
		admin_key = data.AdminKey.ValueString()
	}

	base_url := os.Getenv("OPENAI_BASE_URL")
	if !data.BaseURL.IsNull() {
		base_url = data.BaseURL.ValueString()
	}

	client := openai.NewClient(api_key, admin_key)
	if base_url != "" {
		if parsed, err := url.Parse(base_url); err == nil {
			client.BaseURL = parsed
		}
	}

	// organization_id := os.Getenv("OPENAI_ORGANIZATION_ID")
	// if !data.OrganizationID.IsNull() {
	// 	organization_id = data.OrganizationID.ValueString()
	// }
	// client.OrganizationID = organization_id

	return client
}
