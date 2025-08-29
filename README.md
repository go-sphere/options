# options proto

This is a proto file that defines extra options for RPC methods. It provides method-level option extensions for protobuf services to support custom metadata and routing configurations in Go applications.

## Features

- Method-level option configuration
- Key-value pair metadata system
- Support for different value types (boolean, string, number)
- Extensible extra data mapping
- Integration with code generation plugins
- Designed for use with `protoc-gen-route` and other sphere plugins

## Proto Definition

The core proto file defines:

```protobuf
syntax = "proto3";

package sphere.options;

import "google/protobuf/descriptor.proto";

message KeyValuePair {
  string key = 1;
  oneof value {
    bool flag = 2;
    string text = 3;
    int64 number = 4;
  }
  map<string, string> extra = 5;
}

extend google.protobuf.MethodOptions {
  repeated KeyValuePair options = 501319300;
}
```

## Usage Examples

### Basic Bot Service Configuration

```protobuf
syntax = "proto3";

package bot.v1;

import "sphere/options/options.proto";

message StartRequest {}
message StartResponse {}

service BotService {
  rpc Start(StartRequest) returns (StartResponse) {
    option (sphere.options.options) = {
      key: "bot"
      extra: [
        {
          key: "command"
          value: "start"
        }
      ]
    };
  }
}
```

### Advanced Configuration with Multiple Options

```protobuf
service MenuService {
  rpc UpdateCount(UpdateCountRequest) returns (UpdateCountResponse) {
    option (sphere.options.options) = {
      key: "bot"
      extra: [
        {
          key: "command"
          value: "start"
        },
        {
          key: "callback_query"
          value: "start"
        }
      ]
    };
  }
  
  rpc GetMenu(GetMenuRequest) returns (GetMenuResponse) {
    option (sphere.options.options) = {
      key: "bot"
      extra: [
        {
          key: "command"
          value: "menu"
        },
        {
          key: "description"
          value: "Display main menu"
        },
        {
          key: "admin_only"
          value: "false"
        }
      ]
    };
  }
}
```

### Using Different Value Types

```protobuf
service ConfigService {
  rpc UpdateSettings(UpdateSettingsRequest) returns (UpdateSettingsResponse) {
    option (sphere.options.options) = [
      {
        key: "auth"
        flag: true  // Boolean value
      },
      {
        key: "rate_limit"
        number: 100  // Integer value
      },
      {
        key: "endpoint"
        text: "/api/v1/settings"  // String value
      },
      {
        key: "metadata"
        extra: {
          "timeout": "30s",
          "retry_count": "3"
        }
      }
    ];
  }
}
```

## Integration with buf

Add this dependency to your `buf.yaml`:

```yaml
version: v2
deps:
  - buf.build/go-sphere/options
```

Configure code generation in `buf.gen.yaml`:

```yaml
version: v2
managed:
  enabled: true
  disable:
    - file_option: go_package_prefix
      module: buf.build/go-sphere/options
plugins:
  - local: protoc-gen-go
    out: api
    opt: paths=source_relative
  - local: protoc-gen-route
    out: api
    opt:
      - paths=source_relative
      - options_key=bot
      - request_model=github.com/go-sphere/sphere/social/telegram;Update
      - response_model=github.com/go-sphere/sphere/social/telegram;Message
```

## Generated Code Usage

When used with `protoc-gen-route`, the options generate routing metadata:

```go
// Generated constants for operations
const OperationBotMenuServiceUpdateCount = "/bot.v1.MenuService/UpdateCount"

// Generated extra data from options
var ExtraBotDataMenuServiceUpdateCount = telegram.NewMethodExtraData(map[string]string{
    "callback_query": "start",
    "command":        "start",
})

// Helper function to get extra data by operation
func GetExtraBotDataByMenuServiceOperation(operation string) *telegram.MethodExtraData {
    switch operation {
    case OperationBotMenuServiceUpdateCount:
        return ExtraBotDataMenuServiceUpdateCount
    default:
        return nil
    }
}

// Generated server interface with metadata
type MenuServiceBotServer interface {
    UpdateCount(context.Context, *UpdateCountRequest) (*UpdateCountResponse, error)
}
```

## Option Configuration

### KeyValuePair Structure

- `key`: String identifier for the option
- `value`: One of three types:
  - `flag`: Boolean value for feature flags
  - `text`: String value for textual data
  - `number`: Integer value for numeric data
- `extra`: Map of additional string key-value pairs

### Method Options

- `options`: Repeated KeyValuePair array attached to method definitions
- Supports multiple option entries per method
- Each option can have different value types and extra data

## Common Use Cases

### Bot Commands

Configure Telegram bot command handlers:

```protobuf
rpc HandleStart(StartRequest) returns (StartResponse) {
  option (sphere.options.options) = {
    key: "bot"
    extra: [
      { key: "command", value: "start" },
      { key: "description", value: "Start the bot" }
    ]
  };
}
```

### API Routing

Configure HTTP API endpoints:

```protobuf
rpc GetUser(GetUserRequest) returns (GetUserResponse) {
  option (sphere.options.options) = [
    {
      key: "http"
      extra: [
        { key: "method", value: "GET" },
        { key: "path", value: "/users/{id}" }
      ]
    },
    {
      key: "auth"
      flag: true
    }
  ];
}
```

### Feature Flags

Enable/disable features per method:

```protobuf
rpc AdminOperation(AdminRequest) returns (AdminResponse) {
  option (sphere.options.options) = [
    {
      key: "admin_only"
      flag: true
    },
    {
      key: "rate_limit"
      number: 10
    }
  ];
}
```

## Best Practices

1. **Use consistent key naming**: Stick to a naming convention (e.g., snake_case)
2. **Group related options**: Use the same key for related configuration options
3. **Choose appropriate value types**: Use `flag` for booleans, `text` for strings, `number` for integers
4. **Leverage extra maps**: Use extra data for complex nested configurations
5. **Document option purposes**: Add comments explaining what each option controls
6. **Validate option usage**: Ensure generated code properly handles all defined options

## Integration with Code Generation

The options proto is designed to work seamlessly with sphere's code generation ecosystem:

- **protoc-gen-route**: Generates routing handlers based on option metadata
- **protoc-gen-sphere**: Creates HTTP servers with option-based configuration
- **Custom plugins**: Can leverage options for domain-specific code generation

## Building and Development

Generate the protocol buffer code:

```bash
make generate
```

This will:
1. Generate Go code from proto definitions
2. Format the generated code
3. Run linting checks
4. Apply code formatting and imports