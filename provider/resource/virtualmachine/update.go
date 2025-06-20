package resource_virtualmachine

import (
	"fmt"
	"terraform-provider-vbridge/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Update(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*api.Client)
	vmID := d.Id()

	if d.HasChange("name") {
		err := apiClient.RenameVM(vmID, d.Get("name").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChanges("cores", "memory_size") {
		if d.Get("shutdown_protection").(bool) {
			return fmt.Errorf("shutdown_protection is enabled, turn off protection before updating specs %s", d.Get("name").(string))
		}
		err := apiClient.UpdateVMSpecifications(vmID, d.Get("memory_size").(int), d.Get("cores").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("operating_system_disk_capacity") {
		err := apiClient.ExtendVMDisk(vmID, d.Get("operating_system_disk_guid").(string), d.Get("operating_system_disk_capacity").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChanges("tags", "description", "notes") {
		metadata := BuildMetadataString(d.Get("tags").(map[string]interface{}), d.Get("description").(string), d.Get("notes").(string))
		err := apiClient.SetMetadata(vmID, metadata)
		if err != nil {
			return err
		}
	}

	if d.HasChange("backup_type") {
		err := apiClient.UpdateVMBackupType(vmID, d.Get("backup_type").(string))
		if err != nil {
			return err
		}
	}

	return Read(d, meta)
}
