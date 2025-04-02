package events

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func ReadyEvent(s *discordgo.Session, r *discordgo.Ready) {
	user := s.State.User
	log.Printf("Logged in as %v#%v", user.Username, user.Discriminator)
}
