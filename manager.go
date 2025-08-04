package shards

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Manager facilitates the management of Shards.
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

	// Should state tracking be enabled.
	stateEnabled bool
}

// AddHandler registers an event handler that will be fired anytime the Discord WSAPI event that matches the function fires.
//
// Shouldn't be called after Init or results in undefined behavior.
func (m *Manager) AddHandler(handler interface{}) {
	m.Lock()
	m.handlers = append(m.handlers, handler)
	m.Unlock()
	m.apply(func(shard *Shard) {
		shard.AddHandler(handler)
	})
}

// AddHandlerOnce registers an event handler that will be fired the next time the Discord WSAPI event that matches the function fires.
//
// Calling this method before the Shard is initialized will panic.
func (m *Manager) AddHandlerOnce(handler interface{}) {
	m.apply(func(shard *Shard) {
		shard.AddHandlerOnce(handler)
	})
}

// ApplicationCommandCreate registers an application command for all Shards.
func (m *Manager) ApplicationCommandCreate(guildID string, cmd *discordgo.ApplicationCommand) (errs []error) {
	m.apply(func(shard *Shard) {
		err := shard.ApplicationCommandCreate(guildID, cmd)
		if err != nil {
			errs = append(errs, err)
		}
	})
	return
}

// ApplicationCommandBulkOverwrite registers a series of application commands for all Shards,
// overwriting existing commands.
func (m *Manager) ApplicationCommandBulkOverwrite(guildID string, cmds []*discordgo.ApplicationCommand) (errs []error) {
	m.apply(func(shard *Shard) {
		err := shard.ApplicationCommandBulkOverwrite(guildID, cmds)
		if err != nil {
			errs = append(errs, err)
		}
	})
	return
}

// ApplicationCommandDelete deregisters an application command for all Shards.
func (m *Manager) ApplicationCommandDelete(guildID string, cmd *discordgo.ApplicationCommand) (errs []error) {
	m.apply(func(shard *Shard) {
		err := shard.ApplicationCommandDelete(guildID, cmd)
		if err != nil {
			errs = append(errs, err)
		}
	})
	return
}

// GuildCount returns the amount of guilds that a Manager's Shards are handling.
func (m *Manager) GuildCount() (count int) {
	m.apply(func(shard *Shard) {
		count += shard.GuildCount()
	})
	return
}

// New creates a new Manager with the recommended number of shards.
// After calling New, call Start to begin connecting the shards.
//
// Example:
// mgr := shards.New("Bot TOKEN")
func New(token string) (mgr *Manager, err error) {
	return NewWithConfig(token, DefaultConfig())
}

// NewWithConfig creates a new Manager with the provided configuration.
// After calling NewWithConfig, call Start to begin connecting the shards.
func NewWithConfig(token string, config *Config) (mgr *Manager, err error) {
	if config == nil {
		return nil, fmt.Errorf("error: configuration cannot be nil")
	}

	// Initialize the Manager with provided configuration.
	mgr = &Manager{
		stateEnabled: config.StateEnabled,
		token:        token,
		Intent:       config.Intent,
	}

	// Initialize the gateway.
	mgr.Gateway, err = discordgo.New(token)
	if err != nil {
		return
	}

	resp, err := mgr.Gateway.GatewayBot()
	if err != nil {
		return nil, fmt.Errorf("error: failed to fetch recommended shard count: %v", err)
	}

	if config.ShardCount > 0 && config.ShardCount > resp.Shards {
		mgr.SetShardCount(config.ShardCount)
	} else {
		mgr.SetShardCount(resp.Shards)
	}

	return
}

// SetShardCount sets the shard count.
// The new shard count won't take effect until the Manager is restarted.
func (m *Manager) SetShardCount(count int) {
	if count < 1 {
		return
	}
	m.Lock()
	m.ShardCount = count
	m.Unlock()
}

// RegisterIntent sets the Intent for all Shards' sessions.
//
// Note: Changing the intent will not take effect until the Manager is restarted.
func (m *Manager) RegisterIntent(intent discordgo.Intent) {
	m.Lock()
	m.Intent = intent
	m.Unlock()
}

// SessionForDM returns the proper session for sending and receiving DM's.
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

// Restart restarts the Manager, and rescales if necessary, all with zero downtime.
func (m *Manager) Restart() (nMgr *Manager, err error) {
	// Lock the old Manager for reading.
	m.RLock()

	// Create a new Manager using our token.
	mgr, err := New(m.token)
	if err != nil {
		m.RUnlock()
		return m, err
	}

	// Apply the same handlers.
	for _, handler := range m.handlers {
		mgr.AddHandler(handler)
	}

	// Apply the same intent
	mgr.RegisterIntent(m.Intent)

	// We have no need to lock the old Manager at this point, and
	// starting the new one will take some time.
	m.RUnlock()

	// Start the new Manager so that it can begin handling events.
	err = mgr.Start()
	if err != nil {
		return m, err
	}

	// Shutdown the old Manager. The new Manager is already handling
	// events.
	m.Shutdown()

	return mgr, nil
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
	m.Shards = []*Shard{}
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
		err = shard.Init(m.token, id, m.ShardCount, m.Intent, m.stateEnabled)
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
	m.apply(func(shard *Shard) {
		err = shard.Stop()
		if err != nil {
			return
		}
	})
	return
}

// apply applies a function to all shards.
func (m *Manager) apply(fn func(*Shard)) {
	m.RLock()
	defer m.RUnlock()

	for _, shard := range m.Shards {
		fn(shard)
	}
}
