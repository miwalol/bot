package commands

import (
	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/message"
	"miwa-bot/models"
	"miwa-bot/utils"
)

func Stats(s *discordgo.Session, m *discordgo.MessageCreate) {
	db := utils.DbConnect()
	p := message.NewPrinter(message.MatchLanguage("en"))
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
				Name:   "Users",
				Value:  p.Sprintf("%d", usersCount),
				Inline: true,
			},
			{
				Name:   "Page Views",
				Value:  p.Sprintf("%d", pageViews),
				Inline: true,
			},
			{
				Name:   "Premium Users",
				Value:  p.Sprintf("%d", premiumUsers),
				Inline: true,
			},
		},
	}

	_, _ = s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
}
