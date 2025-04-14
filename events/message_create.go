package events

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"miwa-bot/commands"
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
		commands.Donate(s, m)
	} else if m.Content == "?help" {
		commands.Help(s, m)
	}
}
