package provider

import (
	"context"
	"terraform-provider-vbridge/api"
	virtualmachine_data "terraform-provider-vbridge/data/virtualmachine"
	objectstoragebucket "terraform-provider-vbridge/resource/objectstorage_bucket"
	"terraform-provider-vbridge/resource/virtualmachine"
	additionaldisk "terraform-provider-vbridge/resource/virtualmachine_additionaldisk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Authentication type to use (either 'apiKey' or 'Bearer').",
				ValidateFunc: validation.StringInSlice([]string{"apiKey", "Bearer"}, false),
			},
			"api_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				Description: "API key or token used for authentication. " +
					"If auth_type = 'Bearer', this is your bearer token. " +
					"If auth_type = 'apiKey', this is your API key.",
			},
			"user_email": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"vbridge_virtual_machine":                virtualmachine.Resource(),
			"vbridge_virtual_machine_additionaldisk": additionaldisk.Resource(),
			"vbridge_objectstorage_bucket":           objectstoragebucket.Resource(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"vbridge_virtual_machine": virtualmachine_data.DataSource(),
		},

		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiURL := "https://api.mycloudspace.co.nz"
	authType := d.Get("auth_type").(string)
	apiKey := d.Get("api_key").(string)
	userEmail := d.Get("user_email").(string)

	client, err := api.NewClient(apiURL, authType, apiKey, userEmail)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create API client",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return client, diags
}
