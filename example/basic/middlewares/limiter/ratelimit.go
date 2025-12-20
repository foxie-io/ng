package limiter

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
)

type Limiter struct {
	ng.DefaultControllerInitializer
	ng.DefaultID[Limiter]

	clients map[string]*ClientData
	mutex   sync.RWMutex
	config  *Config
}

type Config struct {
	Limit        int
	Window       time.Duration
	GenerateID   func(ctx context.Context) string
	ErrorHandler func(ctx context.Context) error
}

type ratelimitConfigKey struct{}

type ClientData struct {
	ReqCounts int
	ResetAt   time.Time
	Limit     int
}

// Ensure RateLimiter implements the required interfaces
var _ interface {
	ng.ID
	ng.Guard
} = (*Limiter)(nil)

var defaultConfig = &Config{
	Limit:  100,
	Window: time.Minute,
	GenerateID: func(ctx context.Context) string {
		return "default-client-id"
	},
	ErrorHandler: func(ctx context.Context) error {
		return ng.Respond(ctx, nghttp.NewErrTooManyRequests())
	},
}

// New creates a new RateLimiter instance
func New(config *Config) *Limiter {
	if config == nil {
		panic("config cannot be nil")
	}

	limiter := &Limiter{
		config:  overideOptional(config, defaultConfig),
		clients: make(map[string]*ClientData),
	}

	limiter.StartCleanup(time.Minute)
	return limiter
}

func WithRateLimit(config *Config) ng.Option {
	return ng.WithMetadata(ratelimitConfigKey{}, config)
}

func (rl *Limiter) Allow(ctx context.Context) error {
	config := rl.config

	// Check if there is a route-specific rate limit configuration
	metadata, _ := ng.GetContext(ctx).Route().Core().Metadata(ratelimitConfigKey{})
	if ratelimitConfig, ok := metadata.(*Config); ok {
		config = overideOptional(ratelimitConfig, rl.config)
	}

	// Generate a unique identifier for the client
	id := config.GenerateID(ctx)

	// Lock the mutex to ensure thread-safe access to the clients map
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	client, exists := rl.clients[id]

	// If the client does not exist or their rate limit window has expired, reset their data
	if !exists || now.After(client.ResetAt) {
		client = &ClientData{
			Limit:     config.Limit,
			ReqCounts: 1,
			ResetAt:   now.Add(config.Window),
		}
		rl.clients[id] = client
	} else {
		client.ReqCounts++
	}

	hasReachedLimit := client.ReqCounts > config.Limit
	ng.Store(ctx, client)
	log.Println(client.ReqCounts, "/", client.Limit)

	if hasReachedLimit {
		client.ReqCounts--
		return config.ErrorHandler(ctx)
	}

	return nil
}

// cleanupExpiredClients removes clients whose rate limit window has expired
func (rl *Limiter) cleanupExpiredClients() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	for id, client := range rl.clients {
		if now.After(client.ResetAt) {
			delete(rl.clients, id)
		}
	}
}

func (rl *Limiter) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanupExpiredClients()
		}
	}()
}

func overideOptional(config *Config, defaultConfig *Config) *Config {
	if config == nil {
		return defaultConfig
	}

	if config.Limit == 0 {
		config.Limit = defaultConfig.Limit
	}

	if config.Window == 0 {
		config.Window = defaultConfig.Window
	}

	if config.GenerateID == nil {
		config.GenerateID = defaultConfig.GenerateID
	}

	if config.ErrorHandler == nil {
		config.ErrorHandler = defaultConfig.ErrorHandler
	}

	return config
}
