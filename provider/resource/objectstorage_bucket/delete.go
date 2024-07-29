package objectstoragebucket

import (
	"fmt"
	"terraform-provider-vbridge/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Delete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*api.Client)

	bucketName := d.Id()
	objectstorageTenantId := d.Get("objectstorage_tenant_id").(int)
	canonicalUserId := d.Get("canonical_user_id").(string)

	err := apiClient.DeleteObjectStorageBucket(bucketName, objectstorageTenantId, canonicalUserId)
	if err != nil {
		return fmt.Errorf("error deleting VM: %w", err)
	}

	d.SetId("")

	return nil
}
