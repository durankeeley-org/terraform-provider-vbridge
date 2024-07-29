package additionaldisk

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"capacity": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"storage_profile": {
			Type:     schema.TypeString,
			Required: true,
		},
		"vm_id": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}
