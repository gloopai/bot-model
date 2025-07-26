package model

// StoragePlatform 定义支持的存储类型
type StoragePlatform string

const (
	StoragePlatformOSS StoragePlatform = "oss" // 阿里云 OSS
	StoragePlatformS3  StoragePlatform = "s3"  // AWS S3
)

type StorageClientOptions struct {
	Platform        StoragePlatform `json:"platform"` // "oss" or "s3"
	Endpoint        string          `json:"endpoint"`
	AccessKeyID     string          `json:"accessKeyID"`
	AccessKeySecret string          `json:"accessKeySecret"`
	BucketName      string          `json:"bucketName"`
	Host            string          `json:"host"`
}
