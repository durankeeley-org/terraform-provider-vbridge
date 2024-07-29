package virtualmachine

import (
	"fmt"
	"terraform-provider-vbridge/api"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Delete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*api.Client)
	vmID := d.Id()

	vm, err := apiClient.GetVMDetailedByID(vmID)
	if err != nil {
		return err
	}

	err = apiClient.PowerOffVM(vmID)
	if err != nil {
		return fmt.Errorf("error shutting down VM: %w", err)
	}

	time.Sleep(10 * time.Second)

	err = apiClient.DeleteVM(vmID, vm.Specification.MoRef)
	if err != nil {
		return fmt.Errorf("error deleting VM: %w", err)
	}

	d.SetId("")
	return nil
}
