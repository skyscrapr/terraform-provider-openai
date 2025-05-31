package openai

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/stretchr/testify/assert"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"openai": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func testAccOpenAI(t *testing.T) {
	// Use this to skip tests that might take a long time or cost too much
	if os.Getenv("TF_ACC_OPENAI") == "" {
		t.Skipf("env var TF_ACC_OPENAI not set. Skipping acceptance test due to cost or time")
	}
}

func TestProviderMetadata(t *testing.T) {
	ctx := context.Background()
	raw := New("test")()
	p, ok := raw.(provider.Provider)
	if !ok {
		t.Fatalf("expected provider.Provider, got %T", raw)
	}

	req := provider.MetadataRequest{}
	resp := &provider.MetadataResponse{}

	p.Metadata(ctx, req, resp)

	assert.Equal(t, "openai", resp.TypeName)
	assert.Equal(t, "test", resp.Version)
}

func TestProviderSchema(t *testing.T) {
	ctx := context.Background()
	raw := New("test")()
	p, ok := raw.(provider.Provider)
	if !ok {
		t.Fatalf("expected provider.Provider, got %T", raw)
	}

	req := provider.SchemaRequest{}
	resp := &provider.SchemaResponse{}

	p.Schema(ctx, req, resp)

	assert.NotNil(t, resp.Schema)
	assert.Contains(t, resp.Schema.Attributes, "api_key")
	assert.Contains(t, resp.Schema.Attributes, "admin_key")
	assert.Contains(t, resp.Schema.Attributes, "base_url")
}

func TestConfigureClient_EnvOverride(t *testing.T) {
	// Set env vars
	t.Setenv("OPENAI_BASE_URL", "https://base-url")

	data := OpenAIProviderModel{
		ApiKey:   types.StringNull(),
		AdminKey: types.StringNull(),
		BaseURL:  types.StringNull(),
	}

	client, err := configureClient(data)
	assert.NoError(t, err)
	assert.Equal(t, "https://base-url", client.BaseURL.String())
}

func TestConfigureClient_ConfigOverride(t *testing.T) {
	// Set env vars
	t.Setenv("OPENAI_BASE_URL", "https://base-url")

	data := OpenAIProviderModel{
		ApiKey:   types.StringNull(),
		AdminKey: types.StringNull(),
		BaseURL:  types.StringValue("https://base-url-from-config"),
	}

	client, err := configureClient(data)
	assert.NoError(t, err)
	assert.Equal(t, "https://base-url-from-config", client.BaseURL.String())
}