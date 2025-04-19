package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"miwa-bot/models"
	"miwa-bot/utils"
	"strconv"
	"strings"
)

func formatUsername(username string, displayName *string) string {
	if displayName != nil {
		return *displayName
	}
	return username
}

func getBio(bio *string, typewriterTexts []string, typewriterEnabled bool) string {
	if len(typewriterTexts) > 0 && typewriterEnabled {
		return typewriterTexts[0]
	}
	if bio != nil {
		return *bio
	}
	return "No description set."
}

func getAssetsField(user models.Users) string {
	var lines []string
	if user.AvatarUrl != nil {
		lines = append(lines, fmt.Sprintf("Avatar: [Click here](%s)", *user.AvatarUrl))
	} else {
		lines = append(lines, "Avatar: Not set")
	}
	if user.BackgroundUrl != nil {
		lines = append(lines, fmt.Sprintf("Background: [Click here](https://cdn.miwa.lol/backgrounds/%s)", *user.BackgroundUrl))
	} else {
		lines = append(lines, "Background: Not set")
	}
	if user.CursorUrl != nil {
		lines = append(lines, fmt.Sprintf("Custom cursor: [Click here](%s)", *user.CursorUrl))
	} else {
		lines = append(lines, "Custom cursor: Not set")
	}

	return strings.Join(lines, "\n")
}

func getBackgroundUrl(url *string) string {
	if url != nil {
		return fmt.Sprintf("https://cdn.miwa.lol/backgrounds/%s", *url)
	}
	return ""
}

func getString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func User(s *discordgo.Session, m *discordgo.MessageCreate) {
	id := strings.TrimSpace(m.Content[len("?user"):])
	if len(id) < 1 {
		_, _ = s.ChannelMessageSendReply(m.ChannelID, "❌ **Please provide a username or an alias.**", m.Reference())
		return
	}

	db := utils.DbConnect()
	var user models.Users
	err := db.Model(models.Users{}).Where("handle = ?", id).Or("alias = ?", id).First(&user).Error

	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "❌ **User not found.** Please check the username and try again.")
		return
	}

	color, _ := strconv.ParseUint((*user.AccentColor)[1:], 16, 16)
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    formatUsername(user.Username, user.DisplayName),
			IconURL: getString(user.AvatarUrl),
			URL:     fmt.Sprintf("https://miwa.lol/%s", user.Username),
		},
		Title:       "User Profile",
		Description: getBio(user.Bio, user.TypewriterTexts, user.TypewriterEnabled),
		Color:       int(color),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: getString(user.AvatarUrl),
		},
		Image: &discordgo.MessageEmbedImage{
			URL: getBackgroundUrl(user.BackgroundUrl),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Username",
				Value:  user.Username,
				Inline: true,
			},
			{
				Name: "Account Created",
				// Format string to a Discord timestamp
				Value:  fmt.Sprintf("<t:%d:R>", user.CreatedAt.Unix()),
				Inline: true,
			},
			{
				Name:   "Assets",
				Value:  getAssetsField(user),
				Inline: false,
			},
		},
	}

	_, _ = s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
}
