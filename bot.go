package model

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	Guid                    string   `json:"guid" gorm:"primaryKey;not null"`
	UserId                  int64    `json:"user_id" gorm:"not null"`
	Token                   string   `json:"token" gorm:"type:text;not null;unique"`
	Id                      int64    `json:"id" gorm:"type:integer;not null"`
	IsBot                   bool     `json:"is_bot,omitempty" gorm:"default:false"`
	FirstName               string   `json:"first_name" gorm:"type:text;not null"`
	LastName                string   `json:"last_name,omitempty" gorm:"type:text"`
	UserName                string   `json:"username,omitempty" gorm:"type:text;unique"`
	LanguageCode            string   `json:"language_code,omitempty" gorm:"type:text"`
	CanJoinGroups           bool     `json:"can_join_groups,omitempty" gorm:"default:false"`
	CanReadAllGroupMessages bool     `json:"can_read_all_group_messages,omitempty" gorm:"default:false"`
	SupportsInlineQueries   bool     `json:"supports_inline_queries,omitempty" gorm:"default:false"`
	URL                     string   `json:"url" gorm:"type:text"`
	HasCustomCertificate    bool     `json:"has_custom_certificate" gorm:"default:false"`
	PendingUpdateCount      int      `json:"pending_update_count" gorm:"default:0"`
	IPAddress               string   `json:"ip_address,omitempty" gorm:"type:text"`
	LastErrorDate           int      `json:"last_error_date,omitempty"`
	LastErrorMessage        string   `json:"last_error_message,omitempty" gorm:"type:text"`
	MaxConnections          int      `json:"max_connections,omitempty" gorm:"default:40"`
	AllowedUpdates          []string `json:"allowed_updates,omitempty" gorm:"-"`
	Note                    string   `json:"note,omitempty" gorm:"type:text"`          // 备注信息
	Status                  string   `json:"status" gorm:"type:text;default:'active'"` // 状态: active, inactive, banned
	Createtime              int      `json:"create_time" gorm:"default:0"`
	Updatetime              int      `json:"update_time" gorm:"default:0"`
}

func (t *TelegramBot) TableName() string {
	return "gloop_telegram_bot"
}

// BotMeta 用于存储 Bot 的元信息
type BotAPIService struct {
	Bot  *tgbotapi.BotAPI
	Meta *BotMeta
}

// BotMetaReq 用于获取 Bot 元信息的请求
type GroupMuteMemberReq struct {
	BotGuid  string // Bot 的 GUID
	ChatID   int64  // 聊天 ID
	UserID   int64  // 用户 ID
	Duration int    // 静音时长（秒）
}

// GroupUnMuteMemberReq 用于取消静音群成员的请求
type GroupUnMuteMemberReq struct {
	BotGuid string // Bot 的 GUID
	ChatID  int64  // 聊天 ID
	UserID  int64  // 用户 ID
}
