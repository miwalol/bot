package events

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"miwa-bot/models"
	"miwa-bot/utils"
)

func GuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	// Ignore bots
	if m.User.Bot {
		return
	}

	_ = s.GuildMemberRoleAdd(m.GuildID, m.User.ID, utils.MemberRoleId)

	db := utils.DbConnect()
	var user models.Users
	db.Model(&models.Users{}).Where(models.Users{DiscordId: m.User.ID}).First(&user)

	if user.IsPremium {
		_ = s.GuildMemberRoleAdd(m.GuildID, m.User.ID, utils.PremiumRoleId)
	}
	if user.IsEarlyPremium {
		_ = s.GuildMemberRoleAdd(m.GuildID, m.User.ID, utils.EarlyPremiumRoleId)
	}
	if user.ImageHostAccess {
		_ = s.GuildMemberRoleAdd(m.GuildID, m.User.ID, utils.ImageHostRoleId)
	}
	if user.IsPremiumPlus {
		_ = s.GuildMemberRoleAdd(m.GuildID, m.User.ID, utils.PremiumPlusRoleId)
	}

	msg := fmt.Sprintf(
		"**Welcome to the Miwa.lol Discord server, %s!** Please read the <#%s> & <#%s> channels and enjoy your stay! <@&%s>",
		m.User.Mention(), utils.RulesChannelId, utils.AboutUsChannelId, utils.WelcomeTeamRoleId,
	)
	_, err := s.ChannelMessageSend(utils.ChatChannelId, msg)
	if err != nil {
		log.Fatalf("Failed to send welcome message: %v", err)
	}
}
