# redacted

Generic Go types for YAML that don't leak secrets: they stringify and marshal
as `[REDACTED]` (or just the file path). Built on `go.yaml.in/yaml/v4`.

```sh
go get github.com/mpolitzer/redacted
```

## Types

- `Literal[T]` — inline value. Unmarshals a scalar, marshals `[REDACTED]`.
- `File[T]` — path to a file, read on unmarshal, marshals back to the path. `File[[]byte]` reads raw bytes.
- `Value[T]` — either form: a scalar or `{file: path}`.

```yaml
api_key: "s3cr3t"             # Literal
token: "/etc/app/secret.yml"  # File
foo:
  file: /etc/app/secret.yml   # Value
```

```go
type Conf struct {
    APIKey redacted.Literal[string] `yaml:"api_key"`
    Token  redacted.File[string]    `yaml:"token"`
    Foo    redacted.Value[string]   `yaml:"foo"`
}

c.APIKey.Value() // "s3cr3t"
fmt.Println(c)   // [REDACTED] everywhere
```

Wrong node kinds and unreadable files are unmarshal errors.

```sh
go test ./...
```
