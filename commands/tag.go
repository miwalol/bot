package commands

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

var messages = map[string]string{
	// "trigger": "response",
	// Bug fixed message
	"fixed": "✅ **This bug has been fixed!** Thank you for your report. If you find any other bugs, please report them in the <#1254772076527484978> channel.",
	// Suggestion denied message
	"denied": "❌ **This suggestion has been denied.** After reviewing the suggestion, we have decided that we will not add this feature to Miwa.lol.",
	// Suggestion accepted message
	"accepted": "✅ **This suggestion has been accepted!** We will add this feature to Miwa.lol in the future.",
	// Suggestion pending message. When we are not sure if we'll add the feature or not
	"pending": "⏳ **This suggestion is pending.** We will review this suggestion and get back to you as soon as possible.",
	// Suggestion worked on message.
	"working": "🔨 **We are working on this suggestion!** You'll see it in the next update.",
}
var aliases = map[string][]string{
	// "original": {"alias1", "alias2"},
	"fixed": {"resolved", "fix"},
}

func Tag(s *discordgo.Session, m *discordgo.MessageCreate) {
	tagMsg := strings.TrimSpace(m.Content[len("?tag"):])
	if tagMsg == "" {
		_, _ = s.ChannelMessageSendReply(m.ChannelID, "❌ **Please provide a tag to send.**", m.Reference())
		return
	}

	tag := messages[tagMsg]
	if tag == "" {
		// Check if the tag is an alias
		for original, aliasList := range aliases {
			for _, alias := range aliasList {
				if tagMsg == alias {
					tag = messages[original]
					break
				}
			}
			if tag != "" {
				break
			}
		}

		if tag == "" {
			_, _ = s.ChannelMessageSendReply(m.ChannelID, "❌ **Tag not found.**", m.Reference())
			return
		}
	}

	_ = s.ChannelMessageDelete(m.ChannelID, m.ID)
	_, _ = s.ChannelMessageSend(m.ChannelID, tag)

	canCloseThread := tagMsg == "fixed" || tagMsg == "resolved" || tagMsg == "fix"
	if canCloseThread {
		locked := true
		_, _ = s.ChannelEditComplex(m.ChannelID, &discordgo.ChannelEdit{
			// Lock the thread to prevent further messages
			Locked: &locked,
			// Also archive the thread to put it in the "Older posts" section
			Archived: &locked,
		})
	}
}
