package events

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"miwa-bot/utils"
	"regexp"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from bots
	if m.Author.Bot {
		return
	}

	if m.ChannelID == utils.ProfilesChannelId {
		r := regexp.MustCompile("https://miwa\\.lol/[a-zA-Z0-9_\\-$@]{1,20}")
		match := r.FindStringSubmatch(m.Content)

		// If the Regex is not matching
		if len(match) < 1 {
			_ = s.ChannelMessageDelete(m.ChannelID, m.ID)
		}
	} else if m.ChannelID == utils.BoostsChannelId {
		// Check if the message is a server boost message
		if m.Type != discordgo.MessageTypeUserPremiumGuildSubscription &&
			m.Type != discordgo.MessageTypeUserPremiumGuildSubscriptionTierOne &&
			m.Type != discordgo.MessageTypeUserPremiumGuildSubscriptionTierTwo &&
			m.Type != discordgo.MessageTypeUserPremiumGuildSubscriptionTierThree {
			return
		}
		embed := &discordgo.MessageEmbed{
			Title:       "Thank you for boosting!",
			Description: fmt.Sprintf("Thank you %s for boosting the Miwa Discord server!", m.Author.Mention()),
			Color:       0xF75DDC,
		}

		_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: m.Author.Mention(),
			Embed:   embed,
		})
		if err != nil {
			// Ignore errors
			return
		}
	}

	if m.Content == "?donate" {
		embed := &discordgo.MessageEmbed{
			Title: "Donating",
			Description: fmt.Sprintf(
				"We will be grateful if you can donate to Miwa.lol. Here's the ways to support us by donating:\n" +
					"- [GitHub Sponsors](https://github.com/sponsors/miwalol)\n" +
					"- [PayPal](https://paypal.me/miwalol)",
			),
		}

		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Embed:     embed,
			Reference: m.Reference(),
		})
	} else if m.Content == "?help" {
		embed := &discordgo.MessageEmbed{
			Title: "Help",
			Description: fmt.Sprintf(
				"Here's a list of all the commands:" +
					"- `?donate`: Ways to donate to Miwa.lol" +
					"- `?help`: This help menu",
			),
		}

		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Embed:     embed,
			Reference: m.Reference(),
		})
	}
}
