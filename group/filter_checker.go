package group

import (
	"regexp"
	"strings"
)

// GroupFilterChecker 用于高效检查群组过滤配置
// 用法：checker := NewGroupFilterChecker(config)
// checker.CheckMemberMessage(...), checker.CheckForwardMessage(...), checker.CheckLinkMessage(...), checker.CheckSensitiveWord(...)
type GroupFilterChecker struct {
	Config             *GroupConfig
	sensitiveWordRules []SensitiveWordRule
	domainWhitelistMap map[string]struct{}
}

// NewGroupFilterChecker 构造函数，预编译敏感词正则
func NewGroupFilterChecker(cfg *GroupConfig) *GroupFilterChecker {
	checker := &GroupFilterChecker{Config: cfg}
	// 敏感词直接存规则，检测时用 strings.Contains
	if cfg != nil && cfg.SensitiveWordConfig != nil {
		for _, rule := range cfg.SensitiveWordConfig.Words {
			if rule.Keyword != "" {
				checker.sensitiveWordRules = append(checker.sensitiveWordRules, rule)
			}
		}
	}
	// 域名白名单用 map 优化
	if cfg != nil && cfg.LinkMessageFilter != nil {
		m := make(map[string]struct{}, len(cfg.LinkMessageFilter.DomainWhitelist))
		for _, domain := range cfg.LinkMessageFilter.DomainWhitelist {
			if domain != "" {
				m[domain] = struct{}{}
			}
		}
		checker.domainWhitelistMap = m
	}
	return checker
}

// CheckAll 统一检测入口
// msgType: 消息类型（如 image、bot_command 等）
// content: 消息内容
// isForward: 是否为转发消息
func (c *GroupFilterChecker) CheckAll(msgType, content string, isForward bool) FilterResult {
	// 1. 敏感词检测优先（如需封禁/报警，先处理）
	hit, rule := c.CheckSensitiveWord(content)
	if hit {
		return FilterResult{SensitiveHit: true, SensitiveRule: rule}
	}
	// 2. 成员消息或转发消息类型检测
	var shouldDelete bool
	var reason string
	if isForward {
		shouldDelete, reason = c.CheckForwardMessage(msgType)
	} else {
		shouldDelete, reason = c.CheckMemberMessage(msgType, content)
	}
	if shouldDelete {
		return FilterResult{ShouldDelete: true, Reason: reason}
	}
	// 3. 链接检测
	shouldDelete, reason = c.CheckLinkMessage(content)
	if shouldDelete {
		return FilterResult{ShouldDelete: true, Reason: reason}
	}
	return FilterResult{}
}

// CheckMemberMessage 检查成员消息类型和内容
func (c *GroupFilterChecker) CheckMemberMessage(msgType string, content string) (shouldDelete bool, reason string) {
	f := c.Config.MessageFilter
	if f == nil || f.MemberMessageFilter == nil {
		return false, ""
	}
	m := f.MemberMessageFilter
	switch msgType {
	case "bot_command":
		if m.DeleteBotCommandMsg {
			return true, "delete bot command"
		}
	case "image":
		if m.DeleteImageMsg {
			return true, "delete image"
		}
	case "voice":
		if m.DeleteVoiceMsg {
			return true, "delete voice"
		}
	case "document":
		if m.DeleteDocumentMsg {
			return true, "delete document"
		}
	case "sticker":
		if m.DeleteStickerMsg {
			return true, "delete sticker"
		}
	case "dice":
		if m.DeleteDiceMsg {
			return true, "delete dice"
		}
	}
	return false, ""
}

// CheckForwardMessage 检查转发消息类型
func (c *GroupFilterChecker) CheckForwardMessage(msgType string) (shouldDelete bool, reason string) {
	f := c.Config.MessageFilter
	if f == nil || f.ForwardMessageFilter == nil {
		return false, ""
	}
	m := f.ForwardMessageFilter
	if m.DeleteAllForwardMsg {
		return true, "delete all forward"
	}
	switch msgType {
	case "image":
		if m.DeleteImageMsg {
			return true, "delete image in forward"
		}
	case "animation":
		if m.DeleteAnimationMsg {
			return true, "delete animation in forward"
		}
	case "video":
		if m.DeleteVideoMsg {
			return true, "delete video in forward"
		}
	}
	return false, ""
}

// CheckLinkMessage 检查链接消息
func (c *GroupFilterChecker) CheckLinkMessage(content string) (shouldDelete bool, reason string) {
	f := c.Config.LinkMessageFilter
	if f == nil || !f.DeleteLinkMsg {
		return false, ""
	}
	// 简单检测链接
	linkRe := regexp.MustCompile(`https?://[\w\.-]+`)
	matches := linkRe.FindAllString(content, -1)
	if len(matches) == 0 {
		return false, ""
	}
	// 检查白名单
	for _, link := range matches {
		whitelisted := false
		for domain := range c.domainWhitelistMap {
			if strings.Contains(link, domain) {
				whitelisted = true
				break
			}
		}
		if !whitelisted {
			return true, "delete link: not in whitelist"
		}
	}
	return false, ""
}

// CheckSensitiveWord 检查敏感词，返回第一个命中的规则和处理方式
func (c *GroupFilterChecker) CheckSensitiveWord(content string) (hit bool, rule *SensitiveWordRule) {
	for i := range c.sensitiveWordRules {
		if strings.Contains(content, c.sensitiveWordRules[i].Keyword) {
			return true, &c.sensitiveWordRules[i]
		}
	}
	return false, nil
}
