package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type command struct {
	Name        string
	Description string
}

func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	var description []string
	commands := []command{
		{
			Name:        "help",
			Description: "This command!",
		},
		{
			Name:        "user <username>",
			Description: "Get a user's profile",
		},
		{
			Name:        "me",
			Description: "Get your profile",
		},
		{
			Name:        "donate",
			Description: "Ways to support us by donating",
		},
	}
	for _, cmd := range commands {
		description = append(description, fmt.Sprintf("`?%s`: %s", cmd.Name, cmd.Description))
	}
	embed := &discordgo.MessageEmbed{
		Title:       "Help",
		Description: strings.Join(description[:], "\n"),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Requested by %s", m.Author.Username),
			IconURL: m.Author.AvatarURL(""),
		},
	}

	_, _ = s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
}
