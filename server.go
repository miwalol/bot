package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"miwa-bot/utils"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

func hasNitro(user discordgo.User) bool {
	if user.Banner != "" {
		return true
	}
	return strings.HasPrefix(user.Avatar, "a_")
}

func getMember(s *discordgo.Session, guildID string, userID string, rdb *redis.Client, cacheKey string) (*discordgo.Member, error) {
	member, err := s.State.Member(guildID, userID)

	// If the user is not in the cache
	if err != nil {
		user, err := s.User(userID)
		if err != nil {
			// At this point I give up honestly
			return nil, err
		}

		// Cache that we fetched the user for one week (scroll down to line ≈90 for more info)
		rdb.Set(context.Background(), cacheKey, user.Banner, 7*24*time.Hour)

		// Cache the user in discordgo's state
		member = &discordgo.Member{
			User:    user,
			Avatar:  user.Avatar,
			Banner:  user.Banner,
			Flags:   discordgo.MemberFlags(user.PublicFlags),
			GuildID: guildID,
			Roles:   []string{},
		}
		s.State.MemberAdd(member)
	}
	return member, nil
}

func getClan(clan *discordgo.UserClan) *map[string]any {
	if clan == nil || (clan.Badge == "" && clan.Tag == "") {
		return nil
	}
	return &map[string]any{
		"badge":    clan.Badge,
		"tag":      clan.Tag,
		"guild_id": clan.IdentityGuildId,
	}
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

func getNilOrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func StartServer(s *discordgo.Session) {
	r := gin.Default()

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		ctx := context.Background()
		url, _ := redis.ParseURL(os.Getenv("REDIS_URL"))
		rdb := redis.NewClient(url)
		cacheKey := fmt.Sprintf("discordBanner:%s", id)
		member, err := getMember(s, utils.GuildId, id, rdb, cacheKey)
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

		// Because of Discord API limitations, to check if a user has a banner, we are NEEDED to fetch the user
		// Of course I'm not going to do that every time someone requests the user! so I cache if the user has already been fetched.
		res, err := rdb.Get(ctx, cacheKey).Result()
		if errors.Is(err, redis.Nil) {
			user, err := s.User(id)
			if err != nil {
				log.Printf("Error fetching user: %v", err)
			}

			// Here we just have to update discordgo's members cache.
			if user != nil {
				// Cache the banner for one week
				// In general ppl don't change their banner every hour lol (I hope??), so this should be enough
				rdb.Set(ctx, cacheKey, user.Banner, 7*24*time.Hour)

				// Add the banner to the member object
				member.Banner = user.Banner
				member.User.Banner = user.Banner
				_ = s.State.MemberAdd(member)
			}
		} else if member.User.Banner == "" {
			// If the user is in the cache, we can just get the banner from there
			member.Banner = res
			member.User.Banner = res
		}

		c.JSON(http.StatusOK, gin.H{
			"id":          member.User.ID,
			"username":    member.User.Username,
			"global_name": getNilOrString(member.User.GlobalName),
			"avatar":      getNilOrString(member.User.Avatar),
			// https://support.discord.com/hc/en-us/articles/13410113109911-Avatar-Decorations
			"avatar_decoration": decoration,
			"banner":            getNilOrString(member.Banner),
			"flags":             member.User.PublicFlags,
			// https://support.discord.com/hc/en-us/articles/23187611406999-Guilds-FAQ
			"clan":        getClan(member.User.Clan),
			"has_nitro":   hasNitro(*member.User),
			"is_boosting": isBoosting,
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

		if len(presence.Activities) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"status":   presence.Status,
				"activity": nil,
			})
			return
		}

		// Reversed to get the most recent activity first
		slices.Reverse(presence.Activities)
		activity := presence.Activities[0]

		c.JSON(http.StatusOK, gin.H{
			"status": presence.Status,
			"activity": gin.H{
				"name":           activity.Name,
				"type":           activity.Type,
				"details":        activity.Details,
				"state":          activity.State,
				"emoji":          getStatusEmoji(&activity.Emoji),
				"application_id": activity.ApplicationID,
				"sync_id":        activity.SyncID,
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
