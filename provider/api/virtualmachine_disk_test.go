package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAdditionalDisk(t *testing.T) {
	// Given
	expectedPayload := CreateAdditionalDiskPayload{
		VirtualResourceId: "92582",
		Tier:              "vStorageT1",
		Size:              500,
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle POST /api/virtualresource/AddDisk
		if r.Method == "POST" && r.URL.Path == "/api/virtualresource/AddDisk" {
			w.Header().Set("Content-Type", "application/json")

			var receivedPayload CreateAdditionalDiskPayload
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
	disk := VirtualDisk{
		Capacity:       500,
		StorageProfile: "vStorageT1",
	}
	err := client.CreateAdditionalDisk("92582", disk)

	// Then
	assert.NoError(t, err, "expected no error from CreateAdditionalDisk")
}

func TestCreateAdditionalDiskWithComparison(t *testing.T) {
	// Counter to track the number of GetVMDetailedByID calls
	var getVMDetailedByIDCalls int

	// Given
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle POST /api/virtualresource/AddDisk
		if r.Method == "POST" && r.URL.Path == "/api/virtualresource/AddDisk" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Handle GET /api/VirtualResource/Detailed/{VmId}
		if r.Method == "GET" && r.URL.Path == "/api/VirtualResource/Detailed/12345" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Simulate VM becoming available after the 2nd call
			getVMDetailedByIDCalls++

			firstResponse := map[string]interface{}{
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

			secondResponse := map[string]interface{}{
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
						{
							"moRef":    "8ecc0f6e-633a-40ec-a4c5-4b6463a54305",
							"capacity": 500.0,
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

			var response map[string]interface{}
			if getVMDetailedByIDCalls >= 2 {
				response = secondResponse
			} else {
				response = firstResponse
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	client := testClient(mockServer.URL)

	// When
	disk := VirtualDisk{
		Capacity:       500,
		StorageProfile: "vStorageT1",
	}
	result, err := client.CreateAdditionalDiskWithComparison("12345", disk)

	// Then
	assert.NoError(t, err, "expected no error from CreateAdditionalDiskWithComparison")
	assert.Equal(t, "8ecc0f6e-633a-40ec-a4c5-4b6463a54305", result, "Disk GUID mismatch")
	assert.Equal(t, 2, getVMDetailedByIDCalls, "expected 2 calls to GetVMDetailedByID")
}

func TestFindNewDiskMoRef(t *testing.T) {
	// Given
	initialDisks := []VirtualDisk{
		{MoRef: "disk-100"},
		{MoRef: "disk-101"},
	}

	updatedDisks := []VirtualDisk{
		{MoRef: "disk-100"},
		{MoRef: "disk-101"},
		{MoRef: "disk-102"},
	}

	// When
	result := findNewDiskMoRef(initialDisks, updatedDisks)

	// Then
	assert.Equal(t, "disk-102", result)
}

func TestGetAdditionalDisk(t *testing.T) {
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
	result, err := client.GetVMDisk("12345", "6000C29d-e3d1-85ce-af08-acf6bae05978")

	// Assert
	assert.NoError(t, err)
	disk := VirtualDisk{
		MoRef:          "6000C29d-e3d1-85ce-af08-acf6bae05978",
		Capacity:       100.0,
		StorageProfile: "vStorageT1",
	}
	assert.Equal(t, disk.MoRef, result.MoRef)
	assert.Equal(t, disk.StorageProfile, result.StorageProfile)
	assert.Equal(t, disk.Capacity, result.Capacity)
}
