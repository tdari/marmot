# Creating a Marmot Connection

This guide walks you through creating a connection provider that stores and validates credentials for an external system used by Marmot plugins.

import { CalloutCard } from '@site/src/components/DocCard';

## What is a Connection?

A connection stores the host, port, credentials, and any other authentication configuration needed to reach an external system. Connections are created once in the Marmot UI and reused across many plugin schedules.

## 1. Create the Connection Package

Create a new package under `internal/core/connection/providers`:

```bash
mkdir -p internal/core/connection/providers/hellodb
```

## 2. Implement the Source Interface

Create `source.go` in your connection directory. The connection `Source` interface has a single method:

```go
type Source interface {
    Validate(config map[string]interface{}) error
}
```

Here's a complete example for a fictional "HelloDB" database:

```go
// +marmot:name=HelloDB
// +marmot:description=HelloDB database connection
// +marmot:status=experimental
// +marmot:category=database
package hellodb

//go:generate go run ../../../../docgen/cmd/connection/main.go

import (
    "github.com/marmotdata/marmot/internal/core/connection"
)

// HelloDBConfig defines the configuration for HelloDB connections.
// +marmot:config
type HelloDBConfig struct {
    Host     string `json:"host"     label:"Host"     description:"HelloDB server hostname or IP" validate:"required"           placeholder:"localhost"`
    Port     int    `json:"port"     label:"Port"     description:"HelloDB server port"          default:"5000"  validate:"min=1,max=65535"`
    Database string `json:"database" label:"Database" description:"Database name to connect to"  validate:"required"           placeholder:"mydb"`
    User     string `json:"user"     label:"User"     description:"Username for authentication"  validate:"required"`
    Password string `json:"password" label:"Password" description:"Password for authentication"  sensitive:"true" validate:"required"`
    SSLMode  string `json:"ssl_mode" label:"SSL Mode" description:"SSL connection mode"          default:"disable" validate:"oneof=disable require"`
}

// Example configuration for the connection
// +marmot:example-config
var _ = `
host: localhost
port: 5000
database: mydb
user: admin
password: your-password
ssl_mode: disable
`

// Source implements the connection.Source interface for HelloDB.
type Source struct{}

// Validate validates a raw HelloDB connection config.
func (s *Source) Validate(rawConfig map[string]interface{}) error {
    config, err := connection.UnmarshalConfig[HelloDBConfig](rawConfig)
    if err != nil {
        return err
    }

    // Apply defaults for optional numeric/string fields that may be zero-valued
    if config.Port == 0 {
        config.Port = 5000
    }
    if config.SSLMode == "" {
        config.SSLMode = "disable"
    }

    return connection.ValidateConfig(config)
}

func init() {
    connection.GetRegistry().Register(connection.ConnectionTypeMeta{
        ID:          "hellodb",
        Name:        "HelloDB",
        Description: "HelloDB database connection",
        Icon:        "hellodb",
        Category:    "database",
        ConfigSpec:  connection.GenerateConfigSpec(HelloDBConfig{}),
    }, &Source{})
}
```

## 3. Config Struct Tags

The config struct drives both UI form generation and validation. Every field tag is optional unless noted.

| Tag | Purpose | Example |
|-----|---------|---------|
| `json` | JSON key and serialization | `json:"host"` |
| `label` | Human-readable field name in the UI | `label:"Host"` |
| `description` | Help text shown below the field | `description:"Server hostname"` |
| `validate` | Validation rules (go-playground/validator) | `validate:"required,min=1,max=65535"` |
| `sensitive` | Renders as a password input; encrypted at rest | `sensitive:"true"` |
| `default` | Default value pre-filled in the UI | `default:"5432"` |
| `placeholder` | Ghost text shown in empty input | `placeholder:"localhost"` |
| `show_when` | Show this field only when another field equals a value | `show_when:"connection_string:"` |

### Dropdown fields

Use `validate:"oneof=..."` to render a field as a dropdown select:

```go
SSLMode string `json:"ssl_mode" label:"SSL Mode" validate:"oneof=disable allow require verify-ca verify-full" default:"disable"`
```

### Nested objects

Nest a struct pointer for grouped settings. The UI renders these as a collapsible sub-form:

```go
type TLSConfig struct {
    Enabled  bool   `json:"enabled"   label:"Enabled"   description:"Enable TLS"`
    CertFile string `json:"cert_file" label:"Cert File" description:"Path to client certificate"`
    KeyFile  string `json:"key_file"  label:"Key File"  description:"Path to client key"`
}

type HelloDBConfig struct {
    // ...
    TLS *TLSConfig `json:"tls,omitempty" label:"TLS" description:"TLS configuration"`
}
```

### Conditional fields

Use `show_when` to reveal a field only when another field has a specific value. The format is `field_name:expected_value`. An empty `expected_value` means "show when the field is empty":

```go
// Show AccountName and AccountKey only when connection_string is empty
AccountName string `json:"account_name" show_when:"connection_string:"`
AccountKey  string `json:"account_key"  show_when:"connection_string:" sensitive:"true"`
```

## 4. The `// +marmot:` Annotations

These comments are parsed by the doc generator to produce the connection's documentation page.

| Annotation | Required | Description |
|-----------|---------|-------------|
| `+marmot:name` | Yes | Display name in the UI and docs |
| `+marmot:description` | Yes | One-line description |
| `+marmot:status` | Yes | `stable`, `beta`, or `experimental` |
| `+marmot:category` | Yes | Groups connections in the UI (e.g. `database`, `storage`, `messaging`) |
| `+marmot:config` | Yes | Marks the config struct for docgen extraction |
| `+marmot:example-config` | Recommended | A YAML example shown in the docs; placed on the line before `var _ = \`` |

The `//go:generate` directive on line 7 triggers doc generation when you run `go generate` in the package directory.

## 5. Register the Connection

The `init()` function registers your connection type with the global registry. It runs automatically when the package is imported.

`ConnectionTypeMeta` fields:

| Field | Description |
|-------|-------------|
| `ID` | Unique identifier used internally and in API calls (e.g. `"postgresql"`) |
| `Name` | Human-readable display name |
| `Description` | Short description |
| `Icon` | Icon name. Either an [Iconify](https://icon-sets.iconify.design) identifier or a local SVG filename in `web/docs/static/img/` |
| `Category` | Groups connections in the UI |
| `ConfigSpec` | Auto-generated from your config struct via `connection.GenerateConfigSpec` |

## 6. Import the Package

Import your connection package in `internal/api/v1/server.go` using a blank import to trigger the `init()` registration:

```go
import (
    // existing connection imports ...
    _ "github.com/marmotdata/marmot/internal/core/connection/providers/hellodb"
)
```

Without this import the connection type will not appear in the UI or be available to schedules.

## 7. Generate Documentation

Run `go generate` in your package to produce the connection's documentation page in `web/docs/docs/Connections/`:

```bash
go generate ./internal/core/connection/providers/hellodb/...
```

## 8. Use the Connection in a Plugin

With the connection registered, a plugin can import the config struct and use it to type-assert the merged config at runtime. The plugin's `Discover` method receives a config map that already contains both the connection's credentials and the schedule's discovery settings:

```go
package hellodb

import (
    "context"
    "fmt"

    hellodbc "github.com/marmotdata/marmot/internal/core/connection/providers/hellodb"
    "github.com/marmotdata/marmot/internal/plugin"
)

// Config only contains discovery settings. Credentials come from the linked Connection.
// +marmot:config
type Config struct {
    plugin.BaseConfig `json:",inline"`

    IncludeViews bool `json:"include_views" description:"Discover views in addition to tables" default:"true"`
}

type Source struct {
    config     *Config
    connConfig *hellodbc.HelloDBConfig
}

func (s *Source) Validate(rawConfig plugin.RawPluginConfig) (plugin.RawPluginConfig, error) {
    config, err := plugin.UnmarshalPluginConfig[Config](rawConfig)
    if err != nil {
        return nil, fmt.Errorf("unmarshaling config: %w", err)
    }
    connConfig, err := plugin.UnmarshalPluginConfig[hellodbc.HelloDBConfig](rawConfig)
    if err != nil {
        return nil, fmt.Errorf("unmarshaling connection config: %w", err)
    }
    s.config = config
    s.connConfig = connConfig
    return rawConfig, nil
}

func (s *Source) Discover(ctx context.Context, rawConfig plugin.RawPluginConfig) (*plugin.DiscoveryResult, error) {
    if _, err := s.Validate(rawConfig); err != nil {
        return nil, err
    }

    // Use s.connConfig.Host, s.connConfig.Port, s.connConfig.Password, etc.
    // Use s.config.IncludeViews for discovery behavior.

    return &plugin.DiscoveryResult{}, nil
}
```

The scheduler calls `plugin.MergeConfigs(conn.Config, schedule.Config)` before invoking your plugin, so both sets of fields are present in `rawConfig`. Plugin run config values take precedence over connection config values on key collision.
