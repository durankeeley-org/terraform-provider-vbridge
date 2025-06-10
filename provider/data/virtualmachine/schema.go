package virtualmachine_data

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: Read,
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:      schema.TypeInt,
				Required:  true,
				Sensitive: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"guest_os_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cores": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mo_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hosting_location_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operating_system_disk_guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operating_system_disk_capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"operating_system_disk_storage_profile": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
