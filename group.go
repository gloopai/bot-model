package model

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramGroup struct {
	IDAuto      int64  `gorm:"primaryKey;autoIncrement" json:"id_auto"` // 自增主键
	Guid        string `gorm:"not null" json:"guid"`                    // 群组唯一标识符
	UserId      int64  `gorm:"user_id;default:0" json:"user_id"`        // 用户ID，表示群组所属的用户
	ID          int64  `gorm:"id" json:"id"`                            // Telegram 群组ID
	Title       string `json:"title"`                                   // 群组名称
	Username    string `json:"username"`                                // 群组用户名（可选）
	InviteLink  string `json:"invite_link"`                             // 群组邀请链接（可选）
	BotID       int64  `json:"bot_id" gorm:"default:0"`                 // 绑定的机器人ID
	Type        string `json:"type"`                                    // 群组类型（supergroup, group, channel）
	PhotoURL    string `json:"photo_url"`                               // 群组头像URL
	Description string `json:"description"`                             // 群组描述
	MemberCount int    `json:"member_count"`                            // 群组成员数量
	IsBotAdmin  bool   `json:"is_bot_admin"`                            // 机器人是否为管理员
	BotJoinTime int64  `json:"bot_join_time"`                           // 机器人加入群组的时间
	Language    string `json:"language"`                                // 群组语言
	Permissions string `json:"permissions"`                             // 群组权限（可选，json字符串）
	Config      string `json:"config"`                                  // 群组配置，包含进群选项、消息过滤等
	Status      string `json:"status" gorm:"default:'active'"`          // 群组状态，默认 active
	CreateTime  int    `json:"create_time" gorm:"column:create_time;default:0"`
	UpdateTime  int    `json:"update_time" gorm:"column:update_time;default:0"`
}

func (t *TelegramGroup) TableName() string {
	return "gloop_telegram_group"
}

// UpdateFromChatMemberUpdated 用于将 tgbotapi.ChatMemberUpdated 的信息同步到 TelegramGroup
func (g *TelegramGroup) UpdateFromChatMemberUpdated(update *tgbotapi.ChatMemberUpdated, meta *BotMeta) {
	if update == nil {
		return
	}
	g.UserId = meta.UserId
	g.ID = update.Chat.ID
	g.Title = update.Chat.Title
	g.Username = update.Chat.UserName
	g.Type = update.Chat.Type

	// g.PhotoURL = "" // 需要额外API获取，如：chat, _ := bot.GetChat(tgbotapi.ChatConfig{ChatID: update.Chat.ID}); g.PhotoURL = chat.Photo.BigFileID
	// g.Description = "" // 需要额外API获取，如：chat, _ := bot.GetChat(...); g.Description = chat.Description
	// g.MemberCount = 0 // 需要额外API获取，如：count, _ := bot.GetChatMembersCount(tgbotapi.ChatConfig{ChatID: update.Chat.ID}); g.MemberCount = count
	// 如需获取这些字段，可通过 tgbotapi.GetChat、tgbotapi.GetChatMembersCount 等 API 获取
	if update.NewChatMember.User != nil && update.NewChatMember.User.IsBot {
		g.IsBotAdmin = update.NewChatMember.IsAdministrator()
		g.BotID = update.NewChatMember.User.ID
	}
	// g.BotJoinTime = ... // 可根据事件时间戳设置
	// g.Language = "" // 无法直接获取
	// g.Permissions = "" // 可序列化 update.Chat.Permissions
}

// 群组配置
type GroupConfig struct {
	GroupVerificationConfig GroupVerificationConfig  // 进群验证配置
	MessageFilter           GroupMessageFilter       // 消息过滤配置
	LinkMessageFilter       LinkMessageFilter        // 链接消息过滤
	SensitiveWordConfig     GroupSensitiveWordConfig // 敏感词监控配置
}

// 敏感词监控配置
type GroupSensitiveWordConfig struct {
	Words []GroupSensitiveWordRule // 敏感词规则列表
}

// 敏感词规则
type GroupSensitiveWordRule struct {
	Keyword      string                   // 敏感词
	HandleAction GroupSensitiveWordAction // 处理方式
	BanDuration  int                      // 封禁时长（分钟），仅当处理方式为BanUser时有效
}

// 敏感词处理方式
type GroupSensitiveWordAction int

const (
	GroupSensitiveWordNoAction GroupSensitiveWordAction = iota // 不处理
	GroupSensitiveWordBanUser                                  // 封禁用户
)

/** 消息过滤配置 **/
// 消息过滤配置
type GroupMessageFilter struct {
	MemberMessageFilter  MemberMessageFilter  // 成员消息过滤
	ForwardMessageFilter ForwardMessageFilter // 转发消息过滤
}

// 链接消息过滤
type LinkMessageFilter struct {
	DeleteLinkMsg   bool     // 删除包含链接的消息
	DomainWhitelist []string // 域名白名单
}

// 成员消息过滤
type MemberMessageFilter struct {
	DeleteBotCommandMsg bool // 删除包含机器人指令的消息
	DeleteImageMsg      bool // 删除图片类消息
	DeleteVoiceMsg      bool // 删除语音类消息
	DeleteDocumentMsg   bool // 删除文档类消息
	DeleteStickerMsg    bool // 删除表情和动图类消息
	DeleteDiceMsg       bool // 删除成员投掷的筛子消息
}

// 转发消息过滤
type ForwardMessageFilter struct {
	DeleteImageMsg      bool // 删除包含图片的消息
	DeleteAnimationMsg  bool // 删除包含动画的消息
	DeleteVideoMsg      bool // 删除包含视频的消息
	DeleteAllForwardMsg bool // 删除所有转发消息
}

/*** 进群处理 ***/
// 验证方式枚举
type GroupVerificationMethod int

const (
	GroupVerificationNone     GroupVerificationMethod = iota // 无需验证
	GroupVerificationQuestion                                // 答题验证
	GroupVerificationManual                                  // 管理员手动审核
	// 可扩展更多方式
)

// 验证配置
type GroupVerificationConfig struct {
	WaitTimeSeconds       int                     // 等待认证时间（秒）
	RejoinCooldownSeconds int                     // 重新进入时间（秒）
	Method                GroupVerificationMethod // 验证方式
	PromptMessage         string                  // 提示消息
}

// FilterResult 统一检测结果
type GroupFilterResult struct {
	ShouldDelete  bool                    // 是否应删除消息
	Reason        string                  // 检测结果原因说明
	SensitiveHit  bool                    // 是否命中敏感词
	SensitiveRule *GroupSensitiveWordRule // 命中的敏感词规则（如有）
}
