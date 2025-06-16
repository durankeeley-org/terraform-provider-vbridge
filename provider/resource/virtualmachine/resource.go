package resource_virtualmachine

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema: Schema(),

		Create: Create,
		Read:   Read,
		Update: Update,
		Delete: Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
