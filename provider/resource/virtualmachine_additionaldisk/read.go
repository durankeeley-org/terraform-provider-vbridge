package additionaldisk

import (
	"terraform-provider-vbridge/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Read(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*api.Client)

	vmID := d.Get("vm_id").(string)
	diskID := d.Id()

	vmDisk, err := apiClient.GetVMDisk(vmID, diskID)
	if err != nil {
		return err
	}

	d.Set("capacity", vmDisk.Capacity)
	d.Set("storage_profile", vmDisk.Tier)

	return nil
}
