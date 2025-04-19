package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"log"
	"miwa-bot/utils"
	"net/http"
	"strings"
)

func hasNitro(user discordgo.User) bool {
	if user.Banner != "" {
		return true
	}
	return strings.HasPrefix(user.Avatar, "a_")
}

func getMember(s *discordgo.Session, guildID string, userID string) (*discordgo.Member, error) {
	member, err := s.State.Member(guildID, userID)
	if err != nil {
		// The member is not cached, so we need to fetch it
		member, err = s.GuildMember(guildID, userID)
		if err != nil {
			user, err := s.User(userID)
			if err != nil {
				log.Printf("Error fetching member: %v", err)
				return nil, err
			}
			member = &discordgo.Member{
				User:    user,
				Avatar:  user.Avatar,
				Flags:   discordgo.MemberFlags(user.PublicFlags),
				GuildID: guildID,
				Roles:   []string{},
			}
			s.State.MemberAdd(member)
		}
		s.State.MemberAdd(member)
	}
	return member, nil
}

func getStatusEmoji(emoji *discordgo.Emoji) *map[string]any {
	if emoji == nil || (emoji.ID == "" && emoji.Name == "") {
		return nil
	}
	return &map[string]any{
		"id":       emoji.ID,
		"name":     emoji.Name,
		"animated": emoji.Animated,
	}
}

func StartServer(s *discordgo.Session) {
	r := gin.Default()

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		member, err := getMember(s, utils.GuildId, id)
		if err != nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		isBoosting := false
		if err == nil && member != nil {
			isBoosting = member.PremiumSince != nil
		}

		var decoration *string
		if member.User.AvatarDecorationData != nil {
			decoration = &member.User.AvatarDecorationData.Asset
		}

		c.JSON(http.StatusOK, gin.H{
			"id":                member.User.ID,
			"username":          member.User.Username,
			"global_name":       member.User.GlobalName,
			"avatar":            member.User.Avatar,
			"avatar_decoration": decoration,
			"flags":             member.User.PublicFlags,
			"has_nitro":         hasNitro(*member.User),
			"is_boosting":       isBoosting,
		})
	})
	r.GET("/presence/:id", func(c *gin.Context) {
		id := c.Param("id")
		presence, err := s.State.Presence(utils.GuildId, id)
		if err != nil || presence == nil {
			c.JSON(http.StatusOK, gin.H{
				"status":   "offline",
				"activity": nil,
			})
			return
		}

		activity := presence.Activities[0]
		if len(presence.Activities) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"status":   presence.Status,
				"activity": nil,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": presence.Status,
			"activity": gin.H{
				"name":           activity.Name,
				"type":           activity.Type,
				"details":        activity.Details,
				"state":          activity.State,
				"emoji":          getStatusEmoji(&activity.Emoji),
				"application_id": activity.ApplicationID,
				"assets": gin.H{
					"large_image": activity.Assets.LargeImageID,
					"large_text":  activity.Assets.LargeText,
					"small_image": activity.Assets.SmallImageID,
					"small_text":  activity.Assets.SmallText,
				},
			},
		})
	})

	r.Run(":2007")
}
