package provider

import (
	atuin "atuin-tf/internal/atuin_client"
	"context"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &atuinProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &atuinProvider{
			version: version,
		}
	}
}

// atuinProvider is the provider implementation.
type atuinProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// atuinProviderModel maps provider schema data to a Go type.
type atuinProviderModel struct {
	Host types.String `tfsdk:"host"`
}

// Metadata returns the provider type name.
func (p *atuinProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "atuin"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *atuinProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *atuinProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Atuin client")

	// Retrieve provider data from configuration
	var config atuinProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Atuin API Host",
			"The provider cannot create the Atuin API client as there is an unknown configuration value for the Atuin API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ATUIN_HOST environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("ATUIN_HOST")

	if host == "" {
		host = atuin.API_ENDPOINT
	}

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Atuin API Host",
			"The provider cannot create the Atuin API client as there is a missing or empty value for the atuin API host. "+
				"Set the host value in the configuration or use the atuin_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "atuin_host", host)

	tflog.Debug(ctx, "Creating atuin client")

	// Create a new atuin client using the configuration values
	client := http.Client{}

	// Make the atuin client available during DataSource and Resource
	// type Configure methods.
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Atuin client", map[string]any{"success": true})
}

// Resources defines the resources implemented in the provider.
func (p *atuinProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAtuinUser,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *atuinProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}
