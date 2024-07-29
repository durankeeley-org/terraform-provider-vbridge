package objectstoragebucket

import (
	"terraform-provider-vbridge/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Read(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*api.Client)

	bucketName := d.Id()
	objectstorageTenantId := d.Get("objectstorage_tenant_id").(int)

	bucket, err := apiClient.GetObjectStorageBucket(bucketName, objectstorageTenantId)
	if err != nil {
		return err
	}

	d.Set("bucket_name", bucket)

	return nil
}
