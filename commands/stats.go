package commands

import (
	"github.com/bwmarrin/discordgo"
)

func Stats(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title: "Miwa.lol Stats",
		URL:   "https://miwa.lol/",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Users",
			},
		},
	}

	_, _ = s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
}
