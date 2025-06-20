package resource_virtualmachine

import (
	"fmt"
	"terraform-provider-vbridge/api"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Create(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*api.Client)

	template, templateSet := d.GetOk("template")
	capacity, capacitySet := d.GetOk("operating_system_disk_capacity")

	if !templateSet && !capacitySet {
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
	}

	if capacitySet {
		vm.OperatingSystemDisk.Capacity = capacity.(int)
		vm.OperatingSystemDisk.StorageProfile = d.Get("operating_system_disk_storage_profile").(string)
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

	if err := Read(d, meta); err != nil {
		return err
	}

	if capacitySet {
		userRequestedCapacity := capacity.(int)
		serverCapacity := d.Get("operating_system_disk_capacity").(int)
		diskID := d.Get("operating_system_disk_guid").(string)

		if userRequestedCapacity > serverCapacity {
			err = apiClient.ExtendVMDisk(vmID, diskID, userRequestedCapacity)
			if err != nil {
				return fmt.Errorf("failed to extend disk: %s", err)
			}
		}
	}

	var tags map[string]interface{}
	var description, notes string

	if v, ok := d.GetOk("tags"); ok {
		tags = v.(map[string]interface{})
	} else {
		tags = map[string]interface{}{}
	}

	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
	}

	if v, ok := d.GetOk("notes"); ok {
		notes = v.(string)
	}

	metadata := BuildMetadataString(tags, description, notes)
	err = apiClient.SetMetadata(vmID, metadata)

	time.Sleep(5 * time.Second)
	return Read(d, meta)

}
