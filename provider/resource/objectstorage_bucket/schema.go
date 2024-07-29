package objectstoragebucket

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"bucket_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"objectstorage_tenant_id": {
			Type:      schema.TypeInt,
			Required:  true,
			Sensitive: true,
		},
		"canonical_user_id": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"object_lock": {
			Type:     schema.TypeBool,
			Required: true,
		},
	}
}
