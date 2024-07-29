package virtualmachine

import (
	"fmt"
	"terraform-provider-vbridge/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Create(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*api.Client)

	template, templateSet := d.GetOk("template")
	capacity, capacitySet := d.GetOk("operating_system_disk_capacity")

	if templateSet && capacitySet {
		return fmt.Errorf("`operating_system_disk_capacity` should not be set when `template` is specified")
	} else if !templateSet && !capacitySet {
		return fmt.Errorf("`operating_system_disk_capacity` is required when `template` is not specified")
	}

	vm := api.VirtualMachine{
		ClientId:   d.Get("client_id").(int),
		Name:       d.Get("name").(string),
		Template:   d.Get("template").(string),
		GuestOsId:  d.Get("guest_os_id").(string),
		Cores:      d.Get("cores").(int),
		MemorySize: d.Get("memory_size").(int),
		OperatingSystemDisk: api.VirtualDisk{
			// Capacity:       d.Get("operating_system_disk_capacity").(int),
			StorageProfile: d.Get("operating_system_disk_storage_profile").(string),
		},
		BackupType: d.Get("backup_type").(string),
		HostingLocation: api.HostingLocation{
			Id:             d.Get("hosting_location_id").(string),
			Name:           d.Get("hosting_location_name").(string),
			DefaultNetwork: d.Get("hosting_location_default_network").(string),
		},
		QuoteItem: make(map[string]interface{}), // Initialize with an empty map
	}

	if templateSet {
		vm.Template = template.(string)
	} else {
		vm.OperatingSystemDisk.Capacity = capacity.(int)
	}

	if v, ok := d.GetOk("iso_file"); ok {
		vm.IsoFile = v.(string)
	}

	if v, ok := d.GetOk("quote_item"); ok {
		vm.QuoteItem = v.(map[string]interface{})
	}

	vmID, err := apiClient.CreateVM(vm)
	if err != nil {
		return err
	}

	d.SetId(vmID)
	d.Set("vm_id", vmID)

	return Read(d, meta)
}
