package additionaldisk

import (
	"terraform-provider-vbridge/api"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Delete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*api.Client)

	diskID := d.Id()
	vmID := d.Get("vm_id")

	err := apiClient.DeleteVMDisk(vmID.(string), diskID)
	if err != nil {
		return fmt.Errorf("error deleting VM: %w", err)
	}

	d.SetId("")

	return nil
}
