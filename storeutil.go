package ng

import (
	"context"
	"fmt"
)

// Store store value into context with given key
func Store[T any](ctx context.Context, value T, keys ...PayloadKeyer) {
	key := dynamicKey[T](keys...)
	GetContext(ctx).Storage().Store(key, value)
}

// Load load value from context by given key
func Load[T any](ctx context.Context, keys ...PayloadKeyer) (value T, err error) {
	key := dynamicKey[T](keys...)
	val, loaded := GetContext(ctx).Storage().Load(key)
	if !loaded {
		var zero T
		return zero, fmt.Errorf("not found key: %s", key.PayloadKey())
	}

	expectedType, ok := val.(T)
	if !ok {
		return expectedType, fmt.Errorf("invalid type, expected %T, got %T", val, expectedType)
	}

	return expectedType, nil
}

// Delete delete value from context by given key
func Delete[T any](ctx context.Context, keys ...PayloadKeyer) {
	key := dynamicKey[T](keys...)
	GetContext(ctx).Storage().Delete(key)
}

// LoadOrStore load value from context by given key,
// if not found, store the value into context
func LoadOrStore[T any](ctx context.Context, value T, keys ...PayloadKeyer) (actual T, loaded bool, err error) {
	key := dynamicKey[T](keys...)
	val, loaded := GetContext(ctx).Storage().LoadOrStore(key, value)
	expectedType, ok := val.(T)
	if !ok {
		return expectedType, loaded, fmt.Errorf("invalid type, expected %T, got %T", val, expectedType)
	}
	return expectedType, loaded, nil
}

// MustLoad load value from context by given key,
// panic if not found
func MustLoad[T any](ctx context.Context, keys ...PayloadKeyer) T {
	val, err := Load[T](ctx, keys...)
	if err != nil {
		panic(err)
	}
	return val
}

// MustLoadOrStore load value from context by given key,
// if not found, store the value into context, panic if not found
func MustLoadOrStore[T any](ctx context.Context, value T, keys ...PayloadKeyer) (val T, loaded bool) {
	val, loaded, err := LoadOrStore(ctx, value, keys...)
	if err != nil {
		panic(err)
	}
	return val, loaded
}
