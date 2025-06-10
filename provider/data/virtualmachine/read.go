package virtualmachine_data

import (
	"fmt"
	"terraform-provider-vbridge/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Read(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*api.Client)

	name := d.Get("name").(string)
	clientID := d.Get("client_id").(int)

	vmID, err := apiClient.GetVMByName(name, clientID)
	if err != nil {
		return fmt.Errorf("failed to get VM ID by name %q for client %d: %w", name, clientID, err)
	}

	d.SetId(vmID)

	vm, err := apiClient.GetVMDetailedByID(vmID)
	if err != nil {
		return fmt.Errorf("failed to fetch VM details for ID %s: %w", vmID, err)
	}

	d.Set("client_id", vm.ClientId)
	d.Set("name", vm.Name)
	d.Set("guest_os_id", vm.GuestOsId)
	d.Set("cores", vm.Specification.Cores)
	d.Set("memory_size", vm.Specification.MemoryGb)
	d.Set("mo_ref", vm.Specification.MoRef)
	d.Set("backup_type", vm.Specification.BackupType)
	d.Set("hosting_location_id", vm.Specification.HostingLocationId)
	d.Set("vm_id", vm.Id.String())

	if len(vm.Specification.VirtualDisks) > 0 {
		disk := vm.Specification.VirtualDisks[0]
		d.Set("operating_system_disk_guid", disk.MoRef)
		d.Set("operating_system_disk_capacity", int(disk.Capacity))
		d.Set("operating_system_disk_storage_profile", disk.Tier)
	}

	return nil
}
