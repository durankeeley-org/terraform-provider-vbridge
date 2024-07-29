package api

import "encoding/json"

type VirtualMachine struct {
	ClientId            int                    `json:"clientId"`
	Name                string                 `json:"name"`
	Template            string                 `json:"template"`
	GuestOsId           string                 `json:"guestOsId"`
	Cores               int                    `json:"cores"`
	MemorySize          int                    `json:"memorySize"`
	OperatingSystemDisk VirtualDisk            `json:"operatingSystemDisk"`
	IsoFile             string                 `json:"isoFile,omitempty"`
	QuoteItem           map[string]interface{} `json:"quoteItem"`
	HostingLocation     HostingLocation        `json:"hostingLocation"`
	Id                  json.Number            `json:"id,omitempty"`
	Specification       Specification          `json:"specification"`
	MountedISO          *string                `json:"mountedISO"`
	BackupType          string                 `json:"backupType,omitempty"`
}

type HostingLocation struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	DefaultNetwork string `json:"defaultNetwork"`
}

type Specification struct {
	HealthState       string          `json:"healthState"`
	PowerState        string          `json:"powerState"`
	Cores             int             `json:"cores"`
	Sockets           int             `json:"sockets"`
	MemoryGb          int             `json:"memoryGb"`
	MoRef             string          `json:"moRef"`
	VirtualDisks      []VirtualDisk   `json:"virtualDisks"`
	NetworkDevices    []NetworkDevice `json:"networkDevices"`
	HostingLocationId string          `json:"hostingLocationId"`
	BackupType        string          `json:"backupType"`
}

type NetworkDevice struct {
	Name           string `json:"name"`
	MoRef          string `json:"moRef"`
	NetworkName    string `json:"networkName"`
	MacAddress     string `json:"macAddress"`
	Connected      bool   `json:"connected"`
	StartConnected bool   `json:"startConnected"`
	NetworkId      string `json:"networkId"`
}

type VirtualDisk struct {
	// Detailed returns Capacity as a float
	Capacity int `json:"capacity"`
	// StorageProfile is only used for POST /api/Provisioning/VirtualMachine and only when Template not used
	StorageProfile      string `json:"storageProfile"`
	MoRef               string `json:"moRef,omitempty"`
	Vmfs                string `json:"vmfs,omitempty"`
	SlotInfo            string `json:"slotInfo,omitempty"`
	Tier                string `json:"tier,omitempty"`
	Name                string `json:"name,omitempty"`
	CapacityDescription string `json:"capacityDescription,omitempty"`
}

// Temp Solution
// Convert Capacity from float to int for GetVMDetailedByID
func (vd *VirtualDisk) UnmarshalJSON(data []byte) error {
	type Alias VirtualDisk
	aux := &struct {
		Capacity float64 `json:"capacity"`
		*Alias
	}{
		Alias: (*Alias)(vd),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	vd.Capacity = int(aux.Capacity)
	return nil
}

type PowerOperationPayload struct {
	VirtualResourceId string `json:"VirtualResourceId"`
	Operation         string `json:"Operation"`
}

type DeleteVMOperationPayload struct {
	VirtualResourceId string `json:"VirtualResourceId"`
	CheckToken        string `json:"CheckToken"`
}

type CreateAdditionalDiskPayload struct {
	VirtualResourceId string `json:"virtualResourceId"`
	Tier              string `json:"tier"`
	Size              int    `json:"size"`
}

type ExtendDiskPayload struct {
	VirtualResourceId string `json:"VirtualResourceId"`
	DiskUUID          string `json:"diskUUID"`
	NewSize           int    `json:"newSize"`
	Description       string `json:"description"`
}

type DeleteDiskPayload struct {
	VirtualResourceId string `json:"VirtualResourceId"`
	DiskUUID          string `json:"diskUUID"`
	Description       string `json:"description"`
}
