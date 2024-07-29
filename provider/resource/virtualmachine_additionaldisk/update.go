package additionaldisk

import (
	"terraform-provider-vbridge/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Update(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*api.Client)

	diskID := d.Id()
	vmID := d.Get("vm_id")
	if d.HasChange("capacity") {
		_, newSize := d.GetChange("capacity")
		newDiskSize := newSize.(int)

		err := apiClient.ExtendVMDisk(vmID.(string), diskID, newDiskSize)
		if err != nil {
			return err
		}

		d.Set("capacity", newDiskSize)
	}
	return nil
}
