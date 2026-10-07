package redacted

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.yaml.in/yaml/v4"
)

type fileConf struct {
	Secret File[string] `yaml:"secret"`
}

func TestFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "secret.yml")
	assert.NoError(t, os.WriteFile(path, []byte("secret-value"), 0o600))

	var c fileConf
	assert.NoError(t, yaml.Unmarshal([]byte("secret: "+path), &c))
	assert.Equal(t, path, c.Secret.Path)
	assert.Equal(t, "secret-value", c.Secret.Value())
	assert.Equal(t, "[REDACTED]", c.Secret.String())

	var f File[string]
	assert.NoError(t, yaml.Unmarshal([]byte(path), &f), "top-level")
	assert.Equal(t, "secret-value", f.Value())

	out, err := yaml.Marshal(c)
	assert.NoError(t, err)
	assert.Equal(t, "secret: "+path+"\n", string(out))
	assert.NotContains(t, string(out), "secret-value")
}

func TestFileErrors(t *testing.T) {
	var c fileConf
	assert.Error(t, yaml.Unmarshal([]byte("secret: /nonexistent/secret.yml"), &c), "missing file")

	assert.Error(t, yaml.Unmarshal([]byte("secret: {a: b}"), &c), "non-scalar")
}

func TestFileRawContents(t *testing.T) {
	type rawConf struct {
		Raw File[[]byte] `yaml:"raw"`
	}
	path := filepath.Join(t.TempDir(), "abi.json")
	contents := `contents-are-literal-for-byte-slices`
	assert.NoError(t, os.WriteFile(path, []byte(contents), 0o600))

	var c rawConf
	assert.NoError(t, yaml.Unmarshal([]byte("raw: "+path), &c), "raw contents")
	assert.Equal(t, []byte(contents), c.Raw.Value())
}
