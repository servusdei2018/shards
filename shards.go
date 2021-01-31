package shards

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	// How long to pause between connecting shards.
	TIMELIMIT = time.Second * 5
	// Shards library version. Follows semantic versioning (semver.org).
	VERSION = "1.0.0"
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
