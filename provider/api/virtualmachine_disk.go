package api

import (
	"fmt"
	"time"
)

func (c *Client) CreateAdditionalDisk(vmID string, disk VirtualDisk) error {

	endpoint := "/api/virtualresource/AddDisk"
	payload := CreateAdditionalDiskPayload{
		VirtualResourceId: vmID,
		Tier:              disk.StorageProfile,
		Size:              disk.Capacity,
	}

	resp, err := c.apiRequest("POST", endpoint, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) CreateAdditionalDiskWithComparison(vmID string, disk VirtualDisk) (string, error) {
	initialVM, err := c.GetVMDetailedByID(vmID)
	if err != nil {
		return "", fmt.Errorf("error getting VM details before adding disk: %w", err)
	}
	initialDisks := initialVM.Specification.VirtualDisks

	fmt.Println("Creating additional disk...")
	err = c.CreateAdditionalDisk(vmID, disk)
	if err != nil {
		return "", fmt.Errorf("error creating additional disk: %w", err)
	}

	// Polling parameters
	timeout := 2 * time.Minute
	interval := 5 * time.Second
	startTime := time.Now()

	for {
		if time.Since(startTime) >= timeout {
			return "", fmt.Errorf("timed out waiting for disk to be added to VM %s", vmID)
		}

		updatedVM, err := c.GetVMDetailedByID(vmID)
		if err != nil {
			fmt.Printf("Error getting VM details during polling: %v\n", err)
		} else {
			updatedDisks := updatedVM.Specification.VirtualDisks

			newDiskMoRef := findNewDiskMoRef(initialDisks, updatedDisks)
			if newDiskMoRef != "" {
				fmt.Printf("New disk added with MoRef: %s\n", newDiskMoRef)
				return newDiskMoRef, nil
			}
		}

		time.Sleep(interval)
	}
}

func findNewDiskMoRef(initialDisks, updatedDisks []VirtualDisk) string {
	initialDiskMap := make(map[string]bool)
	for _, disk := range initialDisks {
		initialDiskMap[disk.MoRef] = true
	}

	for _, disk := range updatedDisks {
		if _, found := initialDiskMap[disk.MoRef]; !found {
			return disk.MoRef
		}
	}

	return ""
}

func (c *Client) GetVMDisk(vmID string, diskID string) (*VirtualDisk, error) {
	vm, err := c.GetVMDetailedByID(vmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get VM details: %v", err)
	}

	for _, vmDisk := range vm.Specification.VirtualDisks {
		if vmDisk.MoRef == diskID {
			return &vmDisk, nil
		}
	}

	return nil, fmt.Errorf("disk with MoRef %s not found in VM %s", diskID, vmID)
}

func (c *Client) ExtendVMDisk(vmID string, diskID string, newDiskSize int) error {
	endpoint := "/api/VirtualResource/ExtendDisk"
	payload := ExtendDiskPayload{
		VirtualResourceId: vmID,
		DiskUUID:          diskID,
		NewSize:           newDiskSize,
		Description:       "",
	}

	resp, err := c.apiRequest("POST", endpoint, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) DeleteVMDisk(vmID string, diskID string) error {
	endpoint := "/api/virtualresource/DeleteDisk"
	payload := DeleteDiskPayload{
		VirtualResourceId: vmID,
		DiskUUID:          diskID,
		Description:       "",
	}

	resp, err := c.apiRequest("POST", endpoint, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
