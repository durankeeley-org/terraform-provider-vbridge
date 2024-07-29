package provider

import (
	"context"
	"terraform-provider-vbridge/api"
	objectstoragebucket "terraform-provider-vbridge/resource/objectstorage_bucket"
	"terraform-provider-vbridge/resource/virtualmachine"
	additionaldisk "terraform-provider-vbridge/resource/virtualmachine_additionaldisk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
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
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiURL := d.Get("api_url").(string)
	apiKey := d.Get("api_key").(string)
	userEmail := d.Get("user_email").(string)

	client, err := api.NewClient(apiURL, apiKey, userEmail)
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
