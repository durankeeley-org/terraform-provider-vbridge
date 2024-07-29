package additionaldisk

import (
	"terraform-provider-vbridge/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Create(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*api.Client)

	disk := api.VirtualDisk{
		Capacity:       d.Get("capacity").(int),
		StorageProfile: d.Get("storage_profile").(string),
	}

	vmID := d.Get("vm_id").(string)

	diskID, err := apiClient.CreateAdditionalDiskWithComparison(vmID, disk)
	if err != nil {
		return err
	}

	d.SetId(diskID)

	return Read(d, meta)
}
