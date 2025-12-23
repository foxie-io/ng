package ng

// Storage is used to store key/value pairs in the context

import (
	"fmt"
	"sync"
)

// PayloadKeyer is an interface for defining keys used to store and retrieve payloads in the context.
type PayloadKeyer interface {
	PayloadKey() string
}

// TypeKey is a key type based on generic type T
type TypeKey[T any] struct{}

func (p TypeKey[T]) PayloadKey() string {
	return fmt.Sprintf("%T", p)
}

type PayloadKey string

func (p PayloadKey) PayloadKey() string {
	return "__" + string(p) + "__"
}

type Storage interface {
	// Store store value into context by given key
	Store(key PayloadKeyer, value any)

	// Load load value from context by given key
	Load(key PayloadKeyer) (value any, ok bool)

	// LoadOrStore load value from context by given key,
	// if not found, store the value into context
	LoadOrStore(key PayloadKeyer, value any) (actual any, loaded bool)

	// Delete delete value from context by given key
	Delete(key PayloadKeyer)

	// Clear clear all info stored in context
	Clear()

	// Range iterates over all key/value pairs in the storage.
	Range(f func(key any, value any) bool)
}

var (
	// can be replaced for testing or self-defined storage implementation
	NewDefaultStorage = func() Storage { return NewStorage() }
)

// default store implementation using sync.Map
type storage struct {
	m sync.Map
}

func NewStorage() Storage {
	return &storage{}
}

// Store store value into context with given key
func (s *storage) Store(key PayloadKeyer, value any) {
	s.m.Store(key.PayloadKey(), value)
}

// Load load value from context by given key
func (s *storage) Load(key PayloadKeyer) (value any, ok bool) {
	return s.m.Load(key.PayloadKey())
}

// Delete delete value from context by given key
func (s *storage) Delete(key PayloadKeyer) {
	s.m.Delete(key.PayloadKey())
}

// LoadOrStore load value from context by given key,
// if not found, store the value into context
func (s *storage) LoadOrStore(key PayloadKeyer, value any) (actual any, loaded bool) {
	return s.m.LoadOrStore(key.PayloadKey(), value)
}

// Clear clear all info stored in context
func (s *storage) Clear() {
	s.m.Clear()
}

// Clear clear all info stored in context
func (s *storage) Range(fn func(key any, value any) bool) {
	s.m.Range(fn)
}
