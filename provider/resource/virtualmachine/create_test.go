package virtualmachine_test

import (
    "testing"
    "terraform-provider-vbridge/provider"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Define the providers to be used in tests
var testAccProviders = map[string]*schema.Provider{
    "vbridge": provider.Provider(),
}

// Test configuration
func testAccVirtualMachineConfig_basic() string {
    return `
    provider "vbridge" {
        api_url    = "https://api.example.com"
        api_key    = "my-api-key"
        user_email = "user@example.com"
    }

    resource "vbridge_virtual_machine1" "vm" {
        name  = "test-vm"
        cores = 10
    }`
}

// Test for creating a virtual machine resource
func TestVirtualMachineCreation(t *testing.T) {
	// GIVEN
    resource.Test(t, resource.TestCase{
        Providers: testAccProviders,
        Steps: []resource.TestStep{
            {
				// WHEN
                Config: testAccVirtualMachineConfig_basic(),

				// THEN
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("vbridge_virtual_machine.vm", "name", "test-vm"),
                    resource.TestCheckResourceAttr("vbridge_virtual_machine.vm", "cores", "2"),
                ),
            },
        },
    })
}


