package objectstoragebucket

import (
	"terraform-provider-vbridge/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Create(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*api.Client)

	bucketName := d.Get("bucket_name").(string)
	objectstorageTenantId := d.Get("objectstorage_tenant_id").(int)
	canonicalUserId := d.Get("canonical_user_id").(string)
	objectLock := d.Get("object_lock").(bool)

	err := apiClient.CreateObjectStorageBucket(bucketName, objectstorageTenantId, canonicalUserId, objectLock)
	if err != nil {
		return err
	}

	d.SetId(bucketName)

	return Read(d, meta)
}
