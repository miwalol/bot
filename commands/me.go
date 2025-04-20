package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"miwa-bot/models"
	"miwa-bot/utils"
	"strconv"
)

func Me(s *discordgo.Session, m *discordgo.MessageCreate) {
	db := utils.DbConnect()
	var user models.Users
	err := db.Model(models.Users{}).Where("\"discordId\" = ?", m.Author.ID).First(&user).Error
	if err != nil {
		_, _ = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content:   "❌ **You need to link your Discord account to your Miwa account to use this command.**",
			Reference: m.Reference(),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "Link your Discord account",
							Style: discordgo.LinkButton,
							URL:   "https://help.miwa.lol/how-to/link-discord",
						},
					},
				},
			},
		})
		return
	}

	var views int64
	db.Model(models.PageViews{}).Where("\"userId\" = ?", user.Id).Count(&views)
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
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Views: %v", views),
		},
	}

	_, _ = s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
}
