package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (c *Client) CreateVM(vm VirtualMachine) (string, error) {
	endpoint := "/api/Provisioning/VirtualMachine"
	resp, err := c.apiRequest("POST", endpoint, vm)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error response from API: %s - %s", resp.Status, string(bodyBytes))
	}

	// Polling parameters
	timeout := 2 * time.Minute
	interval := 5 * time.Second
	startTime := time.Now()

	for {
		if time.Since(startTime) >= timeout {
			return "", fmt.Errorf("timed out waiting for VM %s to become available", vm.Name)
		}

		vmID, err := c.GetVMByName(vm.Name, vm.ClientId)
		if err == nil {
			return vmID, nil
		}

		time.Sleep(interval)
	}
}

func (c *Client) GetVMByName(vmName string, clientId int) (string, error) {
	endpoint := fmt.Sprintf("/api/client/virtualresources/%d", clientId)
	resp, err := c.apiRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error response from API: %s - %s", resp.Status, string(bodyBytes))
	}

	var vms []struct {
		Id              int    `json:"id"`
		Name            string `json:"name"`
		HostingLocation string `json:"hostingLocation"`
	}

	err = json.NewDecoder(resp.Body).Decode(&vms)
	if err != nil {
		return "", fmt.Errorf("error decoding JSON response: %w", err)
	}

	for _, vm := range vms {
		if vm.Name == vmName {
			return fmt.Sprintf("%d", vm.Id), nil
		}
	}

	return "", fmt.Errorf("VM with name %s not found", vmName)
}

func (c *Client) GetVMDetailedByID(vmID string) (VirtualMachine, error) {
	endpoint := fmt.Sprintf("/api/VirtualResource/Detailed/%s", vmID)
	resp, err := c.apiRequest("GET", endpoint, nil)
	if err != nil {
		return VirtualMachine{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return VirtualMachine{}, fmt.Errorf("error response from API: %s - %s", resp.Status, string(bodyBytes))
	}

	var temp struct {
		VirtualMachine
		HostingLocation string `json:"hostingLocation"`
	}

	err = json.NewDecoder(resp.Body).Decode(&temp)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("error decoding JSON response: %w", err)
	}

	vm := temp.VirtualMachine
	vm.HostingLocation = HostingLocation{Name: temp.HostingLocation}

	// Define the translation map from backend values to friendly names
	var storageProfileMap = map[string]string{
		"Performance":     "vStorageT1",
		"General Purpose": "vStorageT2",
		"Low Use":         "vStorageT3",
	}

	// Translate the Tier for VirtualDisks
	for i, disk := range vm.Specification.VirtualDisks {
		if friendlyName, ok := storageProfileMap[disk.Tier]; ok {
			vm.Specification.VirtualDisks[i].Tier = friendlyName
		}

		if storageProfile, ok := storageProfileMap[disk.Tier]; ok {
			vm.Specification.VirtualDisks[i].StorageProfile = storageProfile
		}
	}

	return vm, nil
}

func (c *Client) PowerOffVM(vmID string) error {
	endpoint := "/api/virtualresource/poweroperation"
	payload := PowerOperationPayload{
		VirtualResourceId: vmID,
		Operation:         "off",
	}

	vm, err := c.GetVMDetailedByID(vmID)
	if err != nil {
		return fmt.Errorf("failed to get VM details for %s: %w", vmID, err)
	}
	if vm.Specification.PowerState == "off" {
		return nil
	}

	resp, err := c.apiRequest("POST", endpoint, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) DeleteVM(vmID string, moRef string) error {
	endpoint := "/api/virtualresource/delete"
	payload := DeleteVMOperationPayload{
		VirtualResourceId: vmID,
		CheckToken:        moRef,
	}

	resp, err := c.apiRequest("POST", endpoint, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error response from API: %s - %s", resp.Status, string(bodyBytes))
	}

	return nil
}

func (c *Client) SetMetadata(vmID, metadataString string) error {
	endpoint := "/api/virtualresource/annotate"
	payload := MetadataPayload{
		VirtualResourceId: vmID,
		Description:       metadataString,
	}

	resp, err := c.apiRequest("POST", endpoint, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) UpdateVMSpecifications(vmID string, MemorySize int, CoreSize int) error {
	endpoint := "/api/virtualresource/changecomputespecification"

	vm, err := c.GetVMDetailedByID(vmID)
	if err != nil {
		return err
	}

	if vm.Specification.PowerState != "off" {
		err = c.PowerOffVM(vmID)
		if err != nil {
			return fmt.Errorf("failed to power off VM %s before updating specifications: %w", vm.Name, err)
		}
	}

	payload := UpdateVMPayload{
		VirtualResourceId: vmID,
		NewMemorySize:     MemorySize,
		CurrentMemorySize: vm.Specification.MemoryGb,
		NewCoreSize:       CoreSize,
		CurrentCoreSize:   vm.Specification.Cores,
		NewSocketSize:     vm.Specification.Sockets,
		CurrentSocketSize: vm.Specification.Sockets,
	}

	resp, err := c.apiRequest("POST", endpoint, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) RenameVM(vmID, name string) error {
	endpoint := "/api/virtualresource/renameoperation"
	payload := RenameVMPayload{
		VirtualMachineId: vmID,
		NewName:          name,
	}

	resp, err := c.apiRequest("POST", endpoint, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) UpdateVMBackupType(vmID, backupType string) error {
	endpoint := "/api/virtualresource/changebackuptype"
	payload := UpdateBackupTypePayload{
		VirtualResourceId:     vmID,
		Type:                  backupType,
		ApplicationProcessing: true,
	}

	resp, err := c.apiRequest("POST", endpoint, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error response from API: %s - %s", resp.Status, string(bodyBytes))
	}

	return nil
}
