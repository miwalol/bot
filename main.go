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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Cannot load .env file!")
	}

	s, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to start the bot: %v", err)
	}

	s.AddHandler(events.ReadyEvent)
	s.AddHandler(events.GuildMemberAdd)

	err = s.Open()
	if err != nil {
		log.Fatalf("Failed to open the session: %v", err)
	}

	// Ensure the bot session closes properly
	defer s.Close()

	// Listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Received SIGTERM, gracefully shutting down...")
}
