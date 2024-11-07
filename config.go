package shards

import (
	"github.com/bwmarrin/discordgo"
)

// Config stores configuration options for the Manager.
type Config struct {
	Intent       discordgo.Intent
	ShardCount   int
	StateEnabled bool
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		Intent:       discordgo.IntentsAllWithoutPrivileged,
		ShardCount:   1,
		StateEnabled: true,
	}
}
