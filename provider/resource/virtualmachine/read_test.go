package resource_virtualmachine

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
	// GIVEN
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("MOCK REQUEST: %s %s\n", r.Method, r.URL.Path)
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
				"name":      "test-vm-1",
				"guestOS":   "Microsoft Windows Server 2019 (64-bit)",
				"guestOsId": "win2019",
				"clientId":  123,
				"specification": map[string]interface{}{
					"cores":             2,
					"sockets":           1,
					"memoryGb":          4,
					"moRef":             "vm-123",
					"backupType":        "vBackupDisk",
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
				"annotation":      "{Tags: application-service: Core, business-service: Digital Solutions, environment: dev, terraform: true} {description: Jumpbox VM} {notes: This is a test}",
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	resource := schema.TestResourceDataRaw(t, Resource().Schema, map[string]interface{}{
		"name":                                  "test-vm-1",
		"client_id":                             123,
		"template":                              "Windows2022_Standard_30GB",
		"guest_os_id":                           "win2019",
		"cores":                                 2,
		"memory_size":                           4,
		"operating_system_disk_storage_profile": "vStorageT1",
		"iso_file":                              "",
		"quote_item":                            "{}",
		"hosting_location_id":                   "vcchcres",
		"hosting_location_name":                 "Christchurch",
		"hosting_location_default_network":      "CHC-CUST-SDC-WAN",
		"backup_type":                           "vBackupDisk",
	})
	resource.SetId("12345")

	// WHEN
	mockClient := api.MockClient(mockServer.URL)
	err := Read(resource, mockClient)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, "12345", resource.Id())
	assert.Equal(t, 123, resource.Get("client_id"))
	assert.Equal(t, "test-vm-1", resource.Get("name"))
	assert.Equal(t, "win2019", resource.Get("guest_os_id"))
	assert.Equal(t, 2, resource.Get("cores"))
	assert.Equal(t, 4, resource.Get("memory_size"))
	assert.Equal(t, "vm-123", resource.Get("mo_ref"))
	assert.Equal(t, "vBackupDisk", resource.Get("backup_type"))
	assert.Equal(t, "vcchcres", resource.Get("hosting_location_id"))

	debugResourceState(resource)
}
