package models

import "time"

type Users struct {
	Username          string    `gorm:"column:handle"`
	Alias             string    `gorm:"column:alias"`
	DisplayName       *string   `gorm:"column:displayName"`
	Bio               *string   `gorm:"column:bio"`
	DiscordId         string    `gorm:"column:discordId"`
	IsPremium         bool      `gorm:"column:isPremium"`
	IsEarlyPremium    bool      `gorm:"column:isEarlyPremium"`
	IsPremiumPlus     bool      `gorm:"column:isPremiumPlus"`
	ImageHostAccess   bool      `gorm:"column:imageHostAccess"`
	AvatarUrl         *string   `gorm:"column:avatarUrl"`
	BackgroundUrl     *string   `gorm:"column:backgroundUrl"`
	CursorUrl         *string   `gorm:"column:cursorUrl"`
	AccentColor       *string   `gorm:"column:accentColor"`
	TypewriterTexts   []string  `gorm:"column:typewriterTexts;serializer:json"`
	TypewriterEnabled bool      `gorm:"column:typewriterEnabled"`
	CreatedAt         time.Time `gorm:"column:createdAt"`
}
