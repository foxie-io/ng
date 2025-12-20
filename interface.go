package ng

import "fmt"

type PayloadKeyer interface {
	PayloadKey() string
}

type PayloadKey string

type TypeKey[T any] struct{}

func (p TypeKey[T]) PayloadKey() string {
	return fmt.Sprintf("%T", p)
}

func (p PayloadKey) PayloadKey() string {
	return "__" + string(p) + "__"
}
