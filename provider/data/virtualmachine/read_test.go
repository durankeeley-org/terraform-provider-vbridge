package datasource_virtualmachine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"terraform-provider-vbridge/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func debugResourceState(d *schema.ResourceData) {
	fmt.Println("=== Simulated Terraform State ===")
	for k, v := range d.State().Attributes {
		fmt.Printf("%s = %v\n", k, v)
	}
	fmt.Println("=================================")
}

func TestReadVirtualMachine(t *testing.T) {
	// Mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/api/client/virtualresources/123":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]map[string]interface{}{
				{
					"id":              12345,
					"name":            "test-vm-1",
					"hostingLocation": "Christchurch",
				},
			})
		case r.Method == "GET" && r.URL.Path == "/api/VirtualResource/Detailed/12345":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":        12345,
				"name":      "DISKVM0000",
				"guestOS":   "Microsoft Windows Server 2019 (64-bit)",
				"guestOsId": "win2019",
				"clientId":  123,
				"specification": map[string]interface{}{
					"cores":             2,
					"sockets":           1,
					"memoryGb":          4,
					"moRef":             "vm-123",
					"backupType":        "vBackupNone",
					"hostingLocationId": "vcchcres",
					"virtualDisks": []map[string]interface{}{
						{
							"moRef":    "disk-moref-001",
							"capacity": 120.0,
							"tier":     "vStorageT1",
						},
					},
				},
				"hostingLocation": "Christchurch",
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	// Setup Terraform schema and data
	resource := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]interface{}{
		"name":      "test-vm-1",
		"client_id": 123,
	})

	mockClient := api.MockClient(mockServer.URL)

	err := Read(resource, mockClient)
	assert.NoError(t, err)
	debugResourceState(resource)

	assert.Equal(t, "12345", resource.Id())
	assert.Equal(t, 123, resource.Get("client_id"))
	assert.Equal(t, "DISKVM0000", resource.Get("name"))
	assert.Equal(t, "win2019", resource.Get("guest_os_id"))
	assert.Equal(t, 2, resource.Get("cores"))
	assert.Equal(t, 4, resource.Get("memory_size"))
	assert.Equal(t, "vm-123", resource.Get("mo_ref"))
	assert.Equal(t, "vBackupNone", resource.Get("backup_type"))
	assert.Equal(t, "vcchcres", resource.Get("hosting_location_id"))
	assert.Equal(t, "12345", resource.Get("vm_id"))
	assert.Equal(t, "disk-moref-001", resource.Get("operating_system_disk_guid"))
	assert.Equal(t, 120, resource.Get("operating_system_disk_capacity"))
	assert.Equal(t, "vStorageT1", resource.Get("operating_system_disk_storage_profile"))
}
