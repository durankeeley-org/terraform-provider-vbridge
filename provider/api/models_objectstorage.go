package api

type User struct {
	Buckets []Bucket `json:"buckets"`
}

type Bucket struct {
	BucketName string `json:"bucketName"`
}

type CreateObjectStorageBucketPayload struct {
	BucketName string `json:"bucketName"`
	ObjectLock bool   `json:"enableObjectLock"`
}
