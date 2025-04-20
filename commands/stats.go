package commands

import (
	"github.com/bwmarrin/discordgo"
	"miwa-bot/models"
	"miwa-bot/utils"
	"strconv"
)

func Stats(s *discordgo.Session, m *discordgo.MessageCreate) {
	db := utils.DbConnect()
	var usersCount int64
	var premiumUsers int64
	var pageViews int64
	db.Model(models.Users{}).Count(&usersCount)
	db.Model(models.Users{}).Where("\"isPremium\" = ?", true).Or("\"isPremiumPlus\" = ?", true).Count(&premiumUsers)
	db.Model(models.PageViews{}).Count(&pageViews)
	embed := &discordgo.MessageEmbed{
		Title: "Miwa.lol Stats",
		URL:   "https://miwa.lol/",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Users",
				Value: strconv.FormatInt(usersCount, 10),
			},
			{
				Name:  "Page Views",
				Value: strconv.FormatInt(pageViews, 10),
			},
			{
				Name:  "Premium Users",
				Value: strconv.FormatInt(premiumUsers, 10),
			},
		},
	}

	_, _ = s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
}
