package ratelimit

import (
	"context"

	"github.com/foxie-io/ng"
)

type configKey struct{}

func WithConfig(config *Config) ng.Option {
	if config.MetadataKey != "" {
		return ng.WithMetadata(config.MetadataKey, config)
	}

	return ng.WithMetadata(configKey{}, config)
}

func SkipRateLimit() ng.Option {
	return ng.WithSkip(ng.DefaultID[Guard]{})
}

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
