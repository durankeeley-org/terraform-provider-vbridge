package api

import (
	"encoding/json"
	"fmt"
)

func (c *Client) CreateObjectStorageBucket(bucketName string, objectstorageTenantId int, canonicalUserId string, objectLock bool) error {

	endpoint := fmt.Sprintf("/api/ObjectStorage/Tenant/%d/Bucket/%s/add", objectstorageTenantId, canonicalUserId)
	payload := CreateObjectStorageBucketPayload{
		BucketName: bucketName,
		ObjectLock: objectLock,
	}

	resp, err := c.apiRequest("POST", endpoint, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) GetObjectStorageBucket(bucketName string, objectstorageTenantId int) (string, error) {

	endpoint := fmt.Sprintf("/api/ObjectStorage/Detailed/%d", objectstorageTenantId)

	resp, err := c.apiRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response struct {
		Users []struct {
			Buckets []struct {
				BucketName string `json:"bucketName"`
			} `json:"buckets"`
		} `json:"users"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode JSON response: %w", err)
	}

	for _, user := range response.Users {
		for _, bucket := range user.Buckets {
			if bucket.BucketName == bucketName {
				return bucket.BucketName, nil
			}
		}
	}

	return "", fmt.Errorf("bucket %s not found", bucketName)
}

func (c *Client) DeleteObjectStorageBucket(bucketName string, objectstorageTenantId int, canonicalUserId string) error {

	endpoint := fmt.Sprintf("/api/ObjectStorage/Tenant/%d/Bucket/%s/%s/delete", objectstorageTenantId, canonicalUserId, bucketName)

	resp, err := c.apiRequest("POST", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
