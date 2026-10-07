package redacted

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.yaml.in/yaml/v4"
)

type eitherConf struct {
	Foo Value[string] `yaml:"foo"`
}

func TestEitherLiteral(t *testing.T) {
	var c eitherConf
	assert.NoError(t, yaml.Unmarshal([]byte("foo: literal"), &c))
	assert.Equal(t, "literal", c.Foo.Value())
	assert.Equal(t, "[REDACTED]", c.Foo.String())

	out, err := yaml.Marshal(c)
	assert.NoError(t, err)
	assert.Contains(t, string(out), "[REDACTED]")
	assert.NotContains(t, string(out), "literal")
}

func TestEitherFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "secret.yml")
	assert.NoError(t, os.WriteFile(path, []byte("secret-value"), 0o600))

	var c eitherConf
	assert.NoError(t, yaml.Unmarshal([]byte("foo:\n  file: "+path), &c))
	assert.Equal(t, path, c.Foo.File.Path)
	assert.Equal(t, "secret-value", c.Foo.Value())

	out, err := yaml.Marshal(c)
	assert.NoError(t, err)
	assert.Equal(t, "foo: "+path+"\n", string(out))
	assert.NotContains(t, string(out), "secret-value")
}

func TestNewValue(t *testing.T) {
	v := NewValue("secret")
	assert.Equal(t, "secret", v.Value())
	assert.Equal(t, "[REDACTED]", v.String())

	out, err := v.MarshalYAML()
	assert.NoError(t, err)
	assert.Equal(t, "[REDACTED]", out)
}

func TestEitherErrors(t *testing.T) {
	var c eitherConf
	assert.Error(t, yaml.Unmarshal([]byte("foo: [a, b]"), &c), "sequence")
	assert.Error(t, yaml.Unmarshal([]byte("foo:\n  other: x"), &c), "missing file key")
	assert.Error(t, yaml.Unmarshal([]byte("foo:\n  file: {a: b}"), &c), "non-scalar path")
}
