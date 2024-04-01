package shards

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	// TIMELIMIT specifies how long to pause between connecting shards.
	TIMELIMIT = time.Second * 5
	// VERSION specifies the shards module version. Follows semantic versioning (semver.org).
	VERSION = "2.4.0"
)

// A Shard represents a shard.
type Shard struct {
	sync.RWMutex

	// The Discord session handling this Shard.
	Session *discordgo.Session
	// This Shard's ID.
	ID int
	// Total Shard count.
	ShardCount int

	// Event handlers.
	handlers []interface{}
}

// AddHandler registers an event handler for a Shard.
//
// Shouldn't be called after Init or results in undefined behavior.
func (s *Shard) AddHandler(handler interface{}) {
	s.Lock()
	defer s.Unlock()

	s.handlers = append(s.handlers, handler)
}

// ApplicationCommandCreate registers an application command for a Shard.
//
// Shouldn't be called before Initialization.
func (s *Shard) ApplicationCommandCreate(guildID string, cmd *discordgo.ApplicationCommand) error {
	s.Lock()
	defer s.Unlock()

	// Referencing s.Session before Initialization will result in a nil pointer dereference panic.
	if s.Session == nil {
		return fmt.Errorf("error: shard.ApplicationCommandCreate must not be called before shard.Init")
	}

	_, err := s.Session.ApplicationCommandCreate(s.Session.State.User.ID, guildID, cmd)
	return err
}

// ApplicationCommandBulkOverwrite registers a series of application commands for a Shard,
// overwriting existing commands.
//
// Shouldn't be called before Initialization.
func (s *Shard) ApplicationCommandBulkOverwrite(guildID string, cmds []*discordgo.ApplicationCommand) error {
	s.Lock()
	defer s.Unlock()

	// Referencing s.Session before Initialization will result in a nil pointer dereference panic.
	if s.Session == nil {
		return fmt.Errorf("error: shard.ApplicationCommandCreate must not be called before shard.Init")
	}

	_, err := s.Session.ApplicationCommandBulkOverwrite(s.Session.State.User.ID, guildID, cmds)
	return err
}

// ApplicationCommandDelete deregisters an application command for a Shard.
//
// Shouldn't be called before Initialization.
func (s *Shard) ApplicationCommandDelete(guildID string, cmd *discordgo.ApplicationCommand) error {
	s.Lock()
	defer s.Unlock()

	// Referencing s.Session before Initialization will result in a nil pointer dereference panic.
	if s.Session == nil {
		return fmt.Errorf("error: shard.ApplicationCommandCreate must not be called before shard.Init")
	}

	return s.Session.ApplicationCommandDelete(s.Session.State.User.ID, guildID, cmd.ID)
}

// GuildCount returns the amount of guilds that a Shard is handling.
func (s *Shard) GuildCount() (count int) {
	s.RLock()
	defer s.RUnlock()

	if s.Session != nil {
		s.Session.State.RLock()
		count += len(s.Session.State.Guilds)
		s.Session.State.RUnlock()
	}

	return
}

// Init initializes a shard with a bot token, its Shard ID, the total
// amount of shards, and a Discord intent.
func (s *Shard) Init(token string, ID, ShardCount int, intent discordgo.Intent) (err error) {
	s.Lock()
	defer s.Unlock()

	// Apply sharding configuration.
	s.ID = ID
	s.ShardCount = ShardCount

	// Create the session.
	s.Session, err = discordgo.New(token)
	if err != nil {
		return
	}

	// Shard the session.
	s.Session.ShardCount = s.ShardCount
	s.Session.ShardID = s.ID

	// Identify our intent.
	s.Session.Identify.Intents = intent

	// Add handlers to the session.
	for _, handler := range s.handlers {
		s.Session.AddHandler(handler)
	}

	// Connect the shard.
	err = s.Session.Open()

	return
}

// Stop stops a shard.
func (s *Shard) Stop() (err error) {
	s.Lock()
	defer s.Unlock()

	// Close the session.
	err = s.Session.Close()

	return
}
