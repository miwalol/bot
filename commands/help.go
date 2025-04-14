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
	description := []string{"Here's a list of all the commands:"}
	commands := []command{
		{
			Name:        "donate",
			Description: "Ways to support us by donating",
		},
	}
	for _, cmd := range commands {
		description = append(description, fmt.Sprintf("- `%s`: %s", cmd.Name, cmd.Description))
	}
	embed := &discordgo.MessageEmbed{
		Title:       "Help",
		Description: strings.Join(description[:], "\n"),
	}

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embed:     embed,
		Reference: m.Reference(),
	})
}
