package redacted

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.yaml.in/yaml/v4"
)

type literalConf struct {
	Secret Literal[string] `yaml:"api_key"`
}

func TestLiteral(t *testing.T) {
	var c literalConf
	assert.NoError(t, yaml.Unmarshal([]byte("api_key: secret"), &c))
	assert.Equal(t, "secret", c.Secret.Value())
	assert.Equal(t, "[REDACTED]", c.Secret.String())

	var l Literal[string]
	assert.NoError(t, yaml.Unmarshal([]byte("secret"), &l), "top-level")
	assert.Equal(t, "secret", l.Value())

	out, err := yaml.Marshal(c)
	assert.NoError(t, err)
	assert.Contains(t, string(out), "[REDACTED]")
	assert.NotContains(t, string(out), "secret")
}

func TestLiteralError(t *testing.T) {
	var c literalConf
	assert.Error(t, yaml.Unmarshal([]byte("api_key: {a: b}"), &c), "non-scalar")
}
