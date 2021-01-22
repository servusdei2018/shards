package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/servusDei2018/shards"
)

// Global variables.
var (
	Mgr *shards.Manager
)

// Variables used for command line parameters.
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	var err error

	// Create a new shard manager using the provided bot token.
	Mgr, err = shards.New("Bot " + Token)
	if err != nil {
		fmt.Println("[ERROR] Error creating manager,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate
	// events.
	Mgr.AddHandler(messageCreate)
	// Register the onConnect func as a callback for Connect events.
	Mgr.AddHandler(onConnect)

	// In this example, we only care about receiving message events.
	Mgr.RegisterIntent(discordgo.MakeIntent(discordgo.IntentsGuildMessages))

	fmt.Println("[INFO] Starting shard manager...")

	// Start all of our shards and begin listening.
	err = Mgr.Start()
	if err != nil {
		fmt.Println("[ERROR] Error starting manager,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("[SUCCESS] Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Manager.
	Mgr.Shutdown()
}

// This function will be called (due to AddHandler above) every time one
// of our shards connects.
func onConnect(s *discordgo.Session, evt *discordgo.Connect) {
	fmt.Printf("[INFO] Shard #%v connected.\n", s.ShardID)
}

// This function will be called (due to AddHandler above) every time a
// new  message is created on any channel that the authenticated bot has
// access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	// This isn't required in this specific example but it's a good
	// practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	// If the message is "restart" restart the shard manager and rescale
	// if necessary, all with zero down-time.
	if m.Content == "restart" {
		s.ChannelMessageSend(m.ChannelID, "[INFO] Restarting shard manager...")
		fmt.Println("[INFO] Restarting shard manager...")
		if err := Mgr.Restart(); err != nil {
			fmt.Println("[ERROR] Error restarting manager,", err)
		} else {
			s.ChannelMessageSend(m.ChannelID, "[SUCCESS] Manager successfully restarted.")
			fmt.Println("[SUCCESS] Manager successfully restarted.")
		}
	}
}
