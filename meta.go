package model

// BotMeta 保存 bot 的元信息
type BotMeta struct {
	Guid    string    // 唯一标识符
	UserId  int64     // 用户ID
	BotID   string    // Bot ID
	BotName string    // Bot 名称
	Token   string    // Bot Token
	Group   bool      // 是否为群组管理机器人
	Lang    string    // 语言
	Status  BotStatus // 状态
}

// BotStatus 表示 Bot 的状态
type BotStatus string

const (
	BotStatusActive   BotStatus = "active"   // 正常运行
	BotStatusInactive BotStatus = "inactive" // 未激活
	BotStatusPaused   BotStatus = "paused"   // 暂停
	BotStatusError    BotStatus = "error"    // 异常
)
