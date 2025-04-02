package models

type Users struct {
	DiscordId       string `gorm:"column:discordId"`
	IsPremium       bool   `gorm:"column:isPremium"`
	IsEarlyPremium  bool   `gorm:"column:isEarlyPremium"`
	IsPremiumPlus   bool   `gorm:"column:isPremiumPlus"`
	ImageHostAccess bool   `gorm:"column:imageHostAccess"`
}
