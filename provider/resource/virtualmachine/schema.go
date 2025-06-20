package resource_virtualmachine

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Identification and metadata
		"client_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Sensitive:   true,
			Description: "The unique client (org) ID associated with this VM. Sensitive information.",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the virtual machine.",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional description of the virtual machine.",
		},
		"notes": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Additional notes or internal metadata about the VM.",
		},
		"tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Map of key-value tags to assign to the VM.",
		},

		// Template and OS
		"template": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional template name to base the VM on.",
		},
		"guest_os_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the guest operating system.",
		},
		"iso_file": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Path to an ISO file used for installation (optional).",
		},

		// Compute resources
		"cores": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Number of virtual CPU cores assigned to the VM.",
		},
		"memory_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Amount of RAM (in MB) allocated to the VM.",
		},

		// Storage configuration
		"operating_system_disk_guid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The GUID of the operating system disk (auto-generated).",
		},
		"operating_system_disk_capacity": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Capacity of the OS disk in GB. Must be a positive integer.",
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				if v, ok := val.(int); ok && v <= 0 {
					errs = append(errs, fmt.Errorf("%q must be a positive integer, got: %d", key, v))
				}
				return
			},
		},
		"operating_system_disk_storage_profile": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The storage profile to use for the OS disk. Common values include 'vStorageT1', 'vStorageT2' or 'vStorageT3'.",
		},

		// Location and networking
		"hosting_location_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the hosting location.",
		},
		"hosting_location_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the hosting location.",
		},
		"hosting_location_default_network": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The default network assigned to the VM at the hosting location.",
		},

		// Backup and provisioning
		"backup_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The type of backup strategy assigned to the VM. Valid values include 'vBackupNone', 'vBackupDisk'.",
		},
		"quote_item": {
			Type:        schema.TypeMap,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Optional quote-related metadata or pricing details.",
		},

		// System-generated identifiers
		"vm_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The internal identifier of the virtual machine (set after creation).",
		},
		"mo_ref": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The managed object reference for the VM.",
		},

		// Safety / operational flags
		"shutdown_protection": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If true, prevents the VM from being automatically shut down to apply changes (e.g., CPU or memory updates).",
		},
	}
}
