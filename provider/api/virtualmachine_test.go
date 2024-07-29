package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateVM(t *testing.T) {
	// Counter to track the number of GetVMByName calls
	var getVMByNameCalls int

	// Given
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle POST /api/Provisioning/VirtualMachine
		if r.Method == "POST" && r.URL.Path == "/api/Provisioning/VirtualMachine" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(``))
			return
		}

		// Handle GET /api/client/virtualresources/{clientId}
		if r.Method == "GET" && r.URL.Path == "/api/client/virtualresources/123" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Simulate VM becoming available after the 2nd call
			getVMByNameCalls++
			var response []map[string]interface{}
			if getVMByNameCalls >= 2 {
				response = []map[string]interface{}{
					{
						"id":              12345,
						"name":            "test-vm-1",
						"hostingLocation": "Christchurch",
					},
					{
						"id":              12346,
						"name":            "test-vm-2",
						"hostingLocation": "Auckland",
					},
				}
			} else {
				response = []map[string]interface{}{
					{
						"id":              12345,
						"name":            "test-vm-1",
						"hostingLocation": "Christchurch",
					},
				}
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	vm := VirtualMachine{
		ClientId:   123,
		Name:       "test-vm-2",
		Template:   "template-123",
		GuestOsId:  "os-123",
		Cores:      4,
		MemorySize: 8192,
		OperatingSystemDisk: VirtualDisk{
			Capacity: 100,
		},
	}

	client := testClient(mockServer.URL)

	// When
	result, err := client.CreateVM(vm)

	// Then
	assert.NoError(t, err, "expected no error from CreateVM")
	assert.Equal(t, "12346", result, "VM ID mismatch")
	assert.Equal(t, 2, getVMByNameCalls, "expected 2 calls to GetVMByName")
}

func TestGetVMByName(t *testing.T) {
	// Given
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Handle GET /api/client/virtualresources/{clientId}
		if r.Method == "GET" && r.URL.Path == "/api/client/virtualresources/123" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			var response = []map[string]interface{}{
				{
					"id":              12345,
					"name":            "test-vm-1",
					"hostingLocation": "Christchurch",
				},
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	client := testClient(mockServer.URL)

	// When
	result, err := client.GetVMByName("test-vm-1", 123)

	// Then
	assert.NoError(t, err, "expected no error from GetVMByName")
	assert.Equal(t, "12345", result, "VM ID mismatch")
}

func TestGetVMDetailedByID(t *testing.T) {
	// Given
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Handle GET /api/VirtualResource/Detailed/{VmId}
		if r.Method == "GET" && r.URL.Path == "/api/VirtualResource/Detailed/12345" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			var response = map[string]interface{}{
				"clientId": 0,
				"specification": map[string]interface{}{
					"cores":    1,
					"sockets":  4,
					"memoryGb": 4,
					"moRef":    "vm-000",
					"virtualDisks": []map[string]interface{}{
						{
							"moRef":    "6000C29d-e3d1-85ce-af08-acf6bae05978",
							"capacity": 100.0,
							"tier":     "Performance",
						},
					},
					"backupType":        "vBackupNone",
					"hostingLocationId": "vcchcres",
				},
				"id":              12345,
				"name":            "DISKVM0000",
				"hostingLocation": "Christchurch",
				"guestOS":         "Microsoft Windows Server 2019 (64-bit)",
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	client := testClient(mockServer.URL)

	// When
	result, err := client.GetVMDetailedByID("12345")

	// Then
	assert.NoError(t, err, "expected no error from GetVMByName")
	assert.Equal(t, json.Number("12345"), result.Id, "VM ID mismatch")

	assert.Equal(t, "DISKVM0000", result.Name, "VM Name mismatch")
	assert.Equal(t, 1, result.Specification.Cores, "Cores mismatch")
	assert.Equal(t, 4, result.Specification.MemoryGb, "Memory size mismatch")
	assert.Equal(t, "vm-000", result.Specification.MoRef, "MoRef mismatch")

	disk := result.Specification.VirtualDisks[0]
	assert.Equal(t, 100, disk.Capacity, "Disk capacity mismatch")
	assert.Equal(t, "vStorageT1", disk.Tier, "Disk storage profile mismatch")

	assert.Equal(t, "vBackupNone", result.Specification.BackupType, "Backup type mismatch")

	assert.Equal(t, "vcchcres", result.Specification.HostingLocationId, "Hosting Location ID mismatch")
	assert.Equal(t, "Christchurch", result.HostingLocation.Name, "Hosting Location Name mismatch")
}

func TestPowerOffVM(t *testing.T) {
	// Given
	expectedPayload := PowerOperationPayload{
		VirtualResourceId: "7452",
		Operation:         "off",
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle POST /api/virtualresource/poweroperation
		if r.Method == "POST" && r.URL.Path == "/api/virtualresource/poweroperation" {
			w.Header().Set("Content-Type", "application/json")

			var receivedPayload PowerOperationPayload
			err := json.NewDecoder(r.Body).Decode(&receivedPayload)
			if err != nil {
				t.Errorf("Error decoding request body: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			assert.Equal(t, expectedPayload, receivedPayload, "Payload mismatch")
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	client := testClient(mockServer.URL)

	// When
	err := client.PowerOffVM("7452")

	// Then
	assert.NoError(t, err, "expected no error from PowerOffVM")
}

func TestDeleteVM(t *testing.T) {
	// Given
	expectedPayload := DeleteVMOperationPayload{
		VirtualResourceId: "7452",
		CheckToken:        "vm-000",
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle POST /api/virtualresource/poweroperation
		if r.Method == "POST" && r.URL.Path == "/api/virtualresource/delete" {
			w.Header().Set("Content-Type", "application/json")

			var receivedPayload DeleteVMOperationPayload
			err := json.NewDecoder(r.Body).Decode(&receivedPayload)
			if err != nil {
				t.Errorf("Error decoding request body: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			assert.Equal(t, expectedPayload, receivedPayload, "Payload mismatch")
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	client := testClient(mockServer.URL)

	// When
	err := client.DeleteVM("7452", "vm-000")

	// Then
	assert.NoError(t, err, "expected no error from DeleteVM")
}
