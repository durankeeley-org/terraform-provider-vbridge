package virtualmachine

import (
	"terraform-provider-vbridge/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Read(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*api.Client)

	vmID := d.Id()
	vm, err := apiClient.GetVMDetailedByID(vmID)
	if err != nil {
		return err
	}

	d.Set("client_id", vm.ClientId)
	d.Set("name", vm.Name)
	d.Set("guest_os_id", vm.GuestOsId)
	d.Set("cores", vm.Specification.Cores)
	d.Set("memory_size", vm.Specification.MemoryGb)
	d.Set("mo_ref", vm.Specification.MoRef)
	d.Set("operating_system_disk_guid", vm.Specification.VirtualDisks[0].MoRef)
	d.Set("operating_system_disk_capacity", int(vm.Specification.VirtualDisks[0].Capacity))
	d.Set("operating_system_disk_storage_profile", vm.Specification.VirtualDisks[0].Tier)
	d.Set("backup_type", vm.Specification.BackupType)
	d.Set("hosting_location_id", vm.Specification.HostingLocationId)
	d.Set("vm_id", vm.Id.String())

	return nil
}
