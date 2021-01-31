package shards

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Manager struct {
	sync.RWMutex

	// Discord gateway.
	Gateway *discordgo.Session
	// Discord intent.
	Intent discordgo.Intent
	// Shards managed by this Manager.
	Shards []*Shard
	// Total Shard count.
	ShardCount int

	// Event handlers.
	handlers []interface{}
	// Discord bot token.
	token string
}

// AddHandler registers an event handler for all Shards.
func (m *Manager) AddHandler(handler interface{}) {
	m.Lock()
	defer m.Unlock()

	m.handlers = append(m.handlers, handler)

	for _, shard := range m.Shards {
		shard.AddHandler(handler)
	}
}

// GuildCount returns the amount of guilds that a Manager's Shards are
// handling.
func (m *Manager) GuildCount() (count int) {
	m.RLock()
	defer m.RUnlock()

	for _, shard := range m.Shards {
		count += shard.GuildCount()
	}

	return
}

// New creates a new Manager with the recommended number of shards.
// After calling New, call Start to begin connecting the shards.
// You may override ShardCount to use an arbitrary amount of shards.
//
// Example:
// mgr := shards.New("Bot TOKEN")
func New(token string) (mgr *Manager, err error) {
	// Initialize the Manager with provided bot token.
	mgr = &Manager{
		token: token,
	}

	// Initialize the gateway.
	mgr.Gateway, err = discordgo.New(token)

	// Set recommended shard count.
	resp, err := mgr.Gateway.GatewayBot()
	if err != nil {
		return
	}

	if resp.Shards < 1 {
		mgr.ShardCount = 1
	} else {
		mgr.ShardCount = resp.Shards
	}

	return
}

// RegisterIntent sets the Intent for all Shards' sessions.
func (m *Manager) RegisterIntent(intent discordgo.Intent) {
	m.Lock()
	defer m.Unlock()
	m.Intent = intent
}

// SessionForDM returns the proper session for sending and receiving
// DM's.
func (m *Manager) SessionForDM() *discordgo.Session {
	m.RLock()
	defer m.RUnlock()

	// Per Discord documentation, Shard 0 is the only shard which can
	// send and receive DM's.
	//
	// See https://discord.com/developers/docs/topics/gateway#sharding
	return m.Shards[0].Session
}

// SessionForGuild returns the proper session for the specified guild.
func (m *Manager) SessionForGuild(guildID int64) *discordgo.Session {
	m.RLock()
	defer m.RUnlock()

	// Formula to determine which shard handles a guild, from Discord
	// docs.
	//
	// See https://discord.com/developers/docs/topics/gateway#sharding
	return m.Shards[(guildID>>22)%int64(m.ShardCount)].Session
}

// Restart restarts the Manager, and rescales if necessary, all with
// zero downtime.
func (m *Manager) Restart() (err error) {
	// Lock the old Manager for reading.
	m.RLock()

	// Create a new Manager using our token.
	mgr, err := New(m.token)
	if err != nil {
		m.RUnlock()
		return
	}

	// Apply the same handlers.
	for _, handler := range m.handlers {
		mgr.AddHandler(handler)
	}

	// We have no need to lock the old Manager at this point, and
	// starting the new one will take some time.
	m.RUnlock()

	// Start the new Manager so that it can begin handling events.
	err = mgr.Start()
	if err != nil {
		return
	}

	// Shutdown the old Manager. The new Manager is already handling
	// events.
	m.Shutdown()
	// Lock the Manager as we replace the old one with the new one.
	m.Lock()
	defer m.Unlock()
	// Replace the old Manager with the new Manager.
	m = mgr
	// Voila! Zero downtime, because there is always
	// at least one Manager handling events.

	return
}

// Start starts the Manager.
func (m *Manager) Start() (err error) {
	m.Lock()
	defer m.Unlock()

	// Ensure that we have at least one Shard.
	if m.ShardCount < 1 {
		m.ShardCount = 1
	}

	// Initialize Shards.
	for i := 0; i < m.ShardCount; i++ {
		m.Shards = append(m.Shards, &Shard{})
	}

	// Add event handlers to Shards and connect them.
	for id, shard := range m.Shards {
		// Add handlers to this shard.
		for _, handler := range m.handlers {
			shard.AddHandler(handler)
		}
		// Connect shard.
		err = shard.Init(m.token, id, m.ShardCount, m.Intent)
		if err != nil {
			return
		}
		// Ratelimit shard connections.
		if id != len(m.Shards)-1 {
			time.Sleep(TIMELIMIT)
		}
	}

	return
}

// Shutdown gracefully terminates the Manager.
func (m *Manager) Shutdown() (err error) {
	m.Lock()
	defer m.Unlock()

	// Stop all shards.
	for _, shard := range m.Shards {
		if err = shard.Stop(); err != nil {
			return
		}
	}

	return
}
