# Go config

Configure Go application with

## Usage

```go
cl := NewConfigLoader()
if err := cl.LoadConfigFiles(
	"config/default.yaml", 
	"config/env/dev.yaml",
); err := nil {
	panic(err)
}

var cfg Config
if err := cl.Unmarshal(&cfg); err != nil {
	panic(err)
}
```

### Load values from Hashicorp Vault

If you need to load secrets from Vault, you will need a pre-configured Vault client and 
pass it to `WithVaultClient` option.

The default Vault function is registered as `vault` has the following signature:

```go
func vault(path string, key string, defaultValue ...string) (string, error) {}
```

Example:

```yaml
toto:
  foo: {{ vault "secrets/path" "foo_key" }}
  bar: {{ vault "secrets/path" "bar_key" "default_value" }}
```

```go
cl := NewConfigLoader(WithVaultClient(vaultClient))
```

### Change Go template delimiter

In some case, you don't want to use `{{ }}` as delimiters, for example when using Helm. You can change 
the delimiters with `WithDelimiters`.

Example:

```yaml
toto:
  foo: [[ env "FOO_VAL" ]]
```

```go
cl := NewConfigLoader(WithDelimiters("[[", "]]"))
```

### Use custom Go template function

If you need custom template function, you can register one with `WithCustomTemplateFunc`. The function 
must return a string and optionally an error.

Example:

```yaml
toto:
  foo: {{ toUpper "bar" }}
```

```go
cl := NewConfigLoader(WithCustomTemplateFunc("toUpper", func(val string) string {
	return strings.ToUpper(val)
}))
```
