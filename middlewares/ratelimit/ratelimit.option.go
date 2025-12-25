package ratelimit

import (
	"context"

	"github.com/foxie-io/ng"
)

type configKey struct{}

// WithConfig sets the rate limit config in the route metadata.
func WithConfig(config *Config) ng.Option {
	if config.MetadataKey != "" {
		return ng.WithMetadata(config.MetadataKey, config)
	}

	return ng.WithMetadata(configKey{}, config)
}

// SkipRateLimit is used to skip rate limiting for a specific route.
func SkipRateLimit() ng.Option {
	return ng.WithSkip(ng.DefaultID[Guard]{})
}

// GetConfig gets the rate limit config from context metadata.
func GetConfig(ctx context.Context, metadateKey string) (*Config, bool) {
	var (
		key any
	)

	if metadateKey == "" {
		key = configKey{}
	} else {
		key = metadateKey
	}

	data, ok := ng.GetContext(ctx).Route().Core().Metadata(key)
	if !ok {
		return nil, false
	}

	config, ok := data.(*Config)
	if !ok {
		return nil, false
	}

	return config, true
}
