package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"miwa-bot/events"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	_ = godotenv.Load()
	s, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to start the bot: %v", err)
	}

	s.Identify.Presence.Game = discordgo.Activity{
		Type: discordgo.ActivityTypeWatching,
		Name: "miwa.lol | ?help",
	}
	s.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildPresences |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuilds

	s.AddHandler(events.ReadyEvent)
	s.AddHandler(events.GuildMemberAdd)
	s.AddHandler(events.MessageCreate)

	err = s.Open()
	if err != nil {
		log.Fatalf("Failed to open the session: %v", err)
	}

	StartServer(s)

	// Ensure the bot session closes properly
	defer s.Close()
	// Listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Received SIGTERM, gracefully shutting down...")
}
