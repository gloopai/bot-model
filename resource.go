package model

type TelegramResource struct {
	Guid        string `json:"guid" gorm:"primaryKey;not null"`
	UserId      int64  `json:"user_id" gorm:"not null"`
	CatalogId   int64  `json:"catalog_id" gorm:"not null"`
	Title       string `json:"title"`
	Type        string `json:"type" gorm:"not null"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Host        string `json:"host"` // 资源存储的主机头
	Size        int64  `json:"size"` // 资源大小
	Createtime  int    `json:"create_time" gorm:"column:create_time;default:0"`
	Updatetime  int    `json:"update_time" gorm:"column:update_time;default:0"`
}

func (t *TelegramResource) TableName() string {
	return "gloop_telegram_resource"
}

type telegramResourceType struct {
	Photo string `json:"photo"` // 图片
	File  string `json:"file"`  // 文件
	Video string `json:"video"` // 视频
}

var (
	TelegramResourceType = telegramResourceType{
		Photo: "photo",
		File:  "file",
		Video: "video",
	}
)

type resourceFileType struct {
	Jpeg  string `json:"jpeg"` // 图片
	Png   string `json:"png"`  // 文件
	Video string `json:"mp4"`  // 视频
}

var (
	ResourceFileType = resourceFileType{
		Jpeg:  "image/jpeg", // 图片
		Png:   "image/png",  // 文件
		Video: "video/mp4",  // 视频
	}
)

const (
	PlatformCloudflare = "cloudflare"
	PlatformAWS        = "aws"
)

// ResourcePlatform 工厂方法，支持自定义 host
func NewResourcePlatform(platform string, host ...string) ResourcePlatformBase {
	var customHost string
	if len(host) > 0 {
		customHost = host[0]
	}
	switch platform {
	case PlatformCloudflare:
		return &CloudflareResourcePlatform{host: customHost}
	case PlatformAWS:
		return &AWSResourcePlatform{host: customHost}
	default:
		return nil
	}
}

type ResourcePlatformBase interface {
	GetPlatform() string            // 获取平台标识
	GetHost() string                // 获取平台 Host
	GetURL(bucketKey string) string // 根据存储桶 key 生成完整地址
}

type CloudflareResourcePlatform struct {
	host string
}

// GetHost implements ResourcePlatformBase.
func (c *CloudflareResourcePlatform) GetHost() string {
	if c.host != "" {
		return c.host
	}
	return "https://cdn.cloudflare.com"
}

// GetPlatform implements ResourcePlatformBase.
func (c *CloudflareResourcePlatform) GetPlatform() string {
	return PlatformCloudflare
}

// GetURL implements ResourcePlatformBase.
func (c *CloudflareResourcePlatform) GetURL(bucketKey string) string {
	return c.GetHost() + "/" + bucketKey
}

type AWSResourcePlatform struct {
	host string
}

// GetHost implements ResourcePlatformBase.
func (a *AWSResourcePlatform) GetHost() string {
	if a.host != "" {
		return a.host
	}
	return "https://s3.amazonaws.com"
}

// GetPlatform implements ResourcePlatformBase.
func (a *AWSResourcePlatform) GetPlatform() string {
	return PlatformAWS
}

// GetURL implements ResourcePlatformBase.
func (a *AWSResourcePlatform) GetURL(bucketKey string) string {
	return a.GetHost() + "/" + bucketKey
}
