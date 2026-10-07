package redacted

import (
	"errors"
	"os"

	"go.yaml.in/yaml/v4"
)

type File[T any] struct {
	Path  string
	value T
}

// Value returns the file contents.
func (f File[T]) Value() T {
	return f.value
}

func (f *File[T]) String() string {
	return "[REDACTED]"
}

func (f File[T]) MarshalYAML() (any, error) {
	return f.Path, nil
}

func (f *File[T]) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return errors.New("redacted.File: expected scalar")
	}
	f.Path = node.Value

	bytes, err := os.ReadFile(f.Path)
	if err != nil {
		return err
	}
	if b, ok := any(&f.value).(*[]byte); ok {
		*b = bytes // raw contents, regardless of format
		return nil
	}
	return yaml.Unmarshal(bytes, &f.value)
}
