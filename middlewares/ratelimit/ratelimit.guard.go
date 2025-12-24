package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/foxie-io/ng"
)

// Ensure Guard implements the required interfaces
var _ interface {
	ng.ID
	ng.Guard
} = (*Guard)(nil)

type (
	Config struct {
		Limit            int
		Window           time.Duration
		Identifier       func(ctx context.Context) string
		ErrorHandler     func(ctx context.Context) error
		SetHeaderHandler func(ctx context.Context, key, value string)
		MetadataKey      string
	}

	ClientData struct {
		ID        string
		ReqCounts int
		ResetAt   time.Time
		Limit     int
	}

	Guard struct {
		ng.DefaultID[Guard]
		clients map[string]*ClientData
		mutex   sync.RWMutex
		config  *Config
	}
)

// allow replaces the default configuration with user-provided values
var DefaultConfig = &Config{
	Limit:  100,
	Window: time.Minute,
	Identifier: func(ctx context.Context) string {
		return "default-client"
	},
	ErrorHandler: func(ctx context.Context) error {
		return errors.New("rate limit exceeded")
	},
}

// New creates a new Guard instance
func New(config *Config) *Guard {
	if config == nil {
		config = DefaultConfig
	}

	guard := &Guard{
		config:  overrideOptional(config, DefaultConfig),
		clients: make(map[string]*ClientData),
	}

	guard.startCleanup(time.Minute)
	return guard
}

func (g *Guard) Allow(ctx context.Context) error {
	config := g.config

	// Check if there is a route-specific rate limit configuration
	if routeConfig, ok := GetConfig(ctx, config.MetadataKey); ok {
		config = overrideOptional(routeConfig, g.config)
	}

	// Generate a unique identifier for the client
	id := config.Identifier(ctx)

	// Lock the mutex to ensure thread-safe access to the clients map
	g.mutex.Lock()
	defer g.mutex.Unlock()

	now := time.Now()
	client, exists := g.clients[id]

	// If the client does not exist or their rate limit window has expired, reset their data
	if !exists || now.After(client.ResetAt) {
		client = &ClientData{
			ID:        id,
			Limit:     config.Limit,
			ReqCounts: 1,
			ResetAt:   now.Add(config.Window),
		}
		g.clients[id] = client
	} else {
		client.ReqCounts++
	}

	// Set rate limit headers
	if config.SetHeaderHandler != nil {
		remaining := client.Limit - client.ReqCounts
		if remaining < 0 {
			remaining = 0
		}
		config.SetHeaderHandler(ctx, "X-RateLimit-Limit", fmt.Sprintf("%d", client.Limit))
		config.SetHeaderHandler(ctx, "X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		config.SetHeaderHandler(ctx, "X-RateLimit-Reset", client.ResetAt.Format(time.RFC3339))
	}

	hasReachedLimit := client.ReqCounts > config.Limit
	ng.Store(ctx, client)

	if hasReachedLimit {
		client.ReqCounts--
		return config.ErrorHandler(ctx)
	}

	return nil
}

// cleanupExpiredClients removes clients whose rate limit window has expired
func (g *Guard) cleanupExpiredClients() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	now := time.Now()
	for id, client := range g.clients {
		if now.After(client.ResetAt) {
			delete(g.clients, id)
		}
	}
}

func (g *Guard) startCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			g.cleanupExpiredClients()
		}
	}()
}

func overrideOptional(config *Config, defaultConfig *Config) *Config {
	if config == nil {
		return defaultConfig
	}

	if config.Limit == 0 {
		config.Limit = defaultConfig.Limit
	}

	if config.Window == 0 {
		config.Window = defaultConfig.Window
	}

	if config.Identifier == nil {
		config.Identifier = defaultConfig.Identifier
	}

	if config.ErrorHandler == nil {
		config.ErrorHandler = defaultConfig.ErrorHandler
	}

	if config.MetadataKey == "" {
		config.MetadataKey = defaultConfig.MetadataKey
	}

	return config
}
