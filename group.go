package model

// 群组配置
type GroupConfig struct {
	JoinOptions         GroupJoinOptions         `json:"join_options"`          // 进群选项
	MessageFilter       GroupMessageFilter       `json:"message_filter"`        // 消息过滤配置
	LinkMessageFilter   LinkMessageFilter        `json:"link_message_filter"`   // 链接消息过滤
	SensitiveWordConfig GroupSensitiveWordConfig `json:"sensitive_word_config"` // 敏感词监控配置
}

// 敏感词监控配置
type GroupSensitiveWordConfig struct {
	Words []GroupSensitiveWordRule `json:"words"` // 敏感词规则列表
}

// 敏感词规则
type GroupSensitiveWordRule struct {
	Keyword      string                   `json:"keyword"`       // 敏感词
	HandleAction GroupSensitiveWordAction `json:"handle_action"` // 处理方式
	BanDuration  int                      `json:"ban_duration"`  // 封禁时长（秒），仅当处理方式为BanUser时有效
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
	MemberMessageFilter  MemberMessageFilter  `json:"member_message_filter"`  // 成员消息过滤
	ForwardMessageFilter ForwardMessageFilter `json:"forward_message_filter"` // 转发消息过滤
}

// 链接消息过滤
type LinkMessageFilter struct {
	DeleteLinkMsg   bool     `json:"delete_link_msg"`  // 删除包含链接的消息
	DomainWhitelist []string `json:"domain_whitelist"` // 域名白名单
}

// 成员消息过滤
type MemberMessageFilter struct {
	DeleteBotCommandMsg bool `json:"delete_bot_command_msg"` // 删除包含机器人指令的消息
	DeleteImageMsg      bool `json:"delete_image_msg"`       // 删除图片类消息
	DeleteVoiceMsg      bool `json:"delete_voice_msg"`       // 删除语音类消息
	DeleteDocumentMsg   bool `json:"delete_document_msg"`    // 删除文档类消息
	DeleteStickerMsg    bool `json:"delete_sticker_msg"`     // 删除表情和动图类消息
	DeleteDiceMsg       bool `json:"delete_dice_msg"`        // 删除成员投掷的筛子消息
}

// 转发消息过滤
type ForwardMessageFilter struct {
	DeleteImageMsg      bool `json:"delete_image_msg"`       // 删除包含图片的消息
	DeleteAnimationMsg  bool `json:"delete_animation_msg"`   // 删除包含动画的消息
	DeleteVideoMsg      bool `json:"delete_video_msg"`       // 删除包含视频的消息
	DeleteAllForwardMsg bool `json:"delete_all_forward_msg"` // 删除所有转发消息
}

/*** 进群处理 ***/
// 进群方式枚举
type GroupJoinMethod int

const (
	GroupJoinAnyone       GroupJoinMethod = iota // 任何人都可以加入
	GroupJoinInviteOnly                          // 仅限邀请
	GroupJoinVerification                        // 新成员需要通过身份验证
	GroupJoinClosed                              // 关闭进群
)

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

// 进群配置
type GroupJoinOptions struct {
	Method             GroupJoinMethod          // 进群方式
	VerificationConfig *GroupVerificationConfig // 验证配置（仅Method为GroupJoinVerification时有效）
}

// FilterResult 统一检测结果
type GroupFilterResult struct {
	ShouldDelete  bool                    // 是否应删除消息
	Reason        string                  // 检测结果原因说明
	SensitiveHit  bool                    // 是否命中敏感词
	SensitiveRule *GroupSensitiveWordRule // 命中的敏感词规则（如有）
}
