package resource

type TelegramResourceType struct {
	Photo string `json:"photo"` // 图片
	File  string `json:"file"`  // 文件
	Video string `json:"video"` // 视频
}
