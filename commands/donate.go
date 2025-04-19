package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func Donate(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title: "Donating",
		Description: fmt.Sprintf(
			"We will be grateful if you can donate to Miwa.lol. Here's the ways to support us by donating:\n" +
				"- [GitHub Sponsors](https://github.com/sponsors/miwalol)\n" +
				"- [PayPal](https://paypal.me/miwalol)",
		),
	}

	_, _ = s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
}
