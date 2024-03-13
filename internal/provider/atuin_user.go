package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	atuin "terraform-provider-atuin/internal/atuin_client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/tyler-smith/go-bip39"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &AtuinUser{}
	_ resource.ResourceWithImportState = &AtuinUser{}
)

func NewAtuinUser() resource.Resource {
	return &AtuinUser{}
}

// AtuinUser defines the resource implementation.
type AtuinUser struct {
	client *atuin.AtuinClient
}

// AtuinUserModel describes the resource data model.
type AtuinUserModel struct {
	Username  types.String `tfsdk:"username"`
	Password  types.String `tfsdk:"password"`
	Email     types.String `tfsdk:"email"`
	Base64Key types.String `tfsdk:"base64_key"`
	Bip39Key  types.String `tfsdk:"bip39_key"`
}

func (r *AtuinUser) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *AtuinUser) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Atuin user",

		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				MarkdownDescription: "Username of Atuin user",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password of Atuin user",
				Required:            true,
				Sensitive:           true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Email of Atuin user",
				Required:            true,
			},
			"base64_key": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"bip39_key": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func (r *AtuinUser) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*atuin.AtuinClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *AtuinUser) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *AtuinUserModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, data.Username.String())

	_, err := r.client.CreateUser(data.Username.ValueString(), data.Password.ValueString(), data.Email.String())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Atuin user, got error: %s", err))
		return
	}

	// Generate encryption key and add to state
	key, err := atuin.GenerateEncryptionKey()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create encryption key, got error: %s", err))
	}
	data.Base64Key = types.StringValue(key)

	bip39Key, err := atuin.ConvertEncryptionKeyToBip39(data.Base64Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to convert encryption key to bip39, got error: %s", err))
	}
	data.Bip39Key = types.StringValue(bip39Key)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created an Atuin user")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AtuinUser) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AtuinUserModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	_, err := r.client.Login(data.Username.ValueString(), data.Password.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to login user, got error: %s", err))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AtuinUser) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *AtuinUserModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AtuinUser) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AtuinUserModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(data.Username.ValueString(), data.Password.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Atuin user, got error: %s", err))
		return
	}
}

func (r *AtuinUser) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: attr_one,attr_two. Got: %q", req.ID),
		)
		return
	}

	var b64Key, bip39Key string
	if atuin.IsValidBip39(idParts[2]) {
		bip39Key = idParts[2]
		bip39Bytes, err := bip39.EntropyFromMnemonic(bip39Key)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to decode bip39 key, got error: %s", err))
		}
		b64Key = base64.StdEncoding.EncodeToString(bip39Bytes)
	} else {
		b64Key = idParts[2]
		bip39Key, _ = atuin.ConvertEncryptionKeyToBip39(b64Key)
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("username"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("password"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("base64_key"), b64Key)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("bip39_key"), bip39Key)...)
}
