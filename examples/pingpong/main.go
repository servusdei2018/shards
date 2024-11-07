package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/servusdei2018/shards/v2"
)

var (
	mgr   *shards.Manager
	token string
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()

	if token == "" {
		fmt.Println("[ERROR] Bot token is required.")
		os.Exit(1)
	}
}

func main() {
	var err error

	// Create a new shard manager using the provided bot token.
	mgr, err = shards.New("Bot " + token)
	if err != nil {
		fmt.Println("[ERROR] Error creating manager,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate
	// events.
	mgr.AddHandler(messageCreate)
	// Register the onConnect func as a callback for Connect events.
	mgr.AddHandler(onConnect)

	// In this example, we only care about receiving message events.
	mgr.RegisterIntent(discordgo.IntentsGuildMessages)

	fmt.Println("[INFO] Starting shard manager...")

	// Start all of our shards and begin listening.
	err = mgr.Start()
	if err != nil {
		fmt.Println("[ERROR] Error starting manager,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("[SUCCESS] Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Manager.
	fmt.Println("[INFO] Stopping shard manager...")
	mgr.Shutdown()
	fmt.Println("[SUCCESS] Shard manager stopped. Bot is shut down.")
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

	switch m.Content {
	case "ping":
		// If the message is "ping" reply with "Pong!"
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	case "pong":
		// If the message is "pong" reply with "Ping!"
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	case "restart":
		// If the message is "restart" restart the shard manager and rescale
		// if necessary, all with zero down-time.
		var err error
		s.ChannelMessageSend(m.ChannelID, "[INFO] Restarting shard manager...")
		fmt.Println("[INFO] Restarting shard manager...")
		mgr, err = mgr.Restart()
		if err != nil {
			fmt.Println("[ERROR] Error restarting manager,", err)
		} else {
			s.ChannelMessageSend(m.ChannelID, "[SUCCESS] Manager successfully restarted.")
			fmt.Println("[SUCCESS] Manager successfully restarted.")
		}
	}
}
