package redacted

import (
	"errors"

	"go.yaml.in/yaml/v4"
)

type Value[T any] struct {
	Literal Literal[T]
	File    File[T]
}

func NewValue[T any](v T) Value[T] {
	return Value[T]{Literal: NewLiteral(v)}
}

func (e Value[T]) Value() T {
	if e.File.Path != "" {
		return e.File.Value()
	}
	return e.Literal.Value()
}

func (e *Value[T]) String() string {
	return "[REDACTED]"
}

func (e Value[T]) MarshalYAML() (any, error) {
	if e.File.Path != "" {
		return e.File.MarshalYAML()
	}
	return e.Literal.MarshalYAML()
}

func (e *Value[T]) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		return node.Decode(&e.Literal)
	}
	if node.Kind != yaml.MappingNode || len(node.Content) != 2 || node.Content[0].Value != "file" {
		return errors.New("redacted.Value: expected scalar or {file: path}")
	}
	return node.Content[1].Decode(&e.File)
}
