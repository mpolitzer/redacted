package redacted

import (
	"errors"

	"go.yaml.in/yaml/v4"
)

type Literal[T any] struct {
	value T
}

// NewLiteral returns a Literal holding v.
func NewLiteral[T any](v T) Literal[T] {
	return Literal[T]{value: v}
}

// Value returns the held value.
func (v Literal[T]) Value() T {
	return v.value
}

func (v *Literal[T]) String() string {
	return "[REDACTED]"
}

func (r Literal[T]) MarshalYAML() (any, error) {
	return "[REDACTED]", nil
}

func (r *Literal[T]) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return errors.New("redacted.Literal: expected scalar")
	}
	return node.Decode(&r.value)
}
