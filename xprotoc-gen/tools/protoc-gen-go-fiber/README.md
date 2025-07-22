# protoc-gen-go-fiber

🇬🇧 [EN](README.md) · 🇷🇺 [RU](README_ru.md)

**protoc-gen-go-fiber** is a plugin for `protoc` or `buf` that automatically generates HTTP routes
for [Fiber](https://github.com/gofiber/fiber) based on gRPC services and `google.api.http` annotations.

## Features

* HTTP handlers generation for Fiber from gRPC services.
* Supports methods: `GET`, `POST`, `PUT`, `PATCH`, `DELETE`.
* Request validation via `protoc-gen-validate`.
* Interceptor support for error handling and custom logic.
* Automatic conversion of HTTP headers into gRPC metadata.
* **Important:** only the request body is accepted.

## Installation

```bash
go install github.com/petara94/protoc-gen-go-fiber@latest
```

Or for Go 1.24+:

```bash
go get -tool github.com/petara94/protoc-gen-go-fiber@latest
```

## Usage

Example using `buf`:

1. Ensure your `.proto` files include `google.api.http` annotations:

```protobuf
import "google/api/annotations.proto";
import "validate/validate.proto";

option go_package = "gen/go/greeterpb;greeterpb";

service GreeterService {
  rpc SayHello(HelloRequest) returns (HelloResponse) {
    option (google.api.http) = {
      post: "api/v1/hello"
      body: "*"
    };
  }

  // New POST method example
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {
    option (google.api.http) = {
      post: "/api/v1/user"
      body: "*"
    };
  }
}

message HelloRequest {
  string name = 1 [(validate.rules).string.min_len = 3];
  string email = 2 [(validate.rules).string.email = true];
}

// New request/response messages for CreateUser
message CreateUserRequest {
  string username = 1 [(validate.rules).string.min_len = 3];
  string email = 2 [(validate.rules).string.email = true];
}

message CreateUserResponse {
  string id = 1;
  string username = 2;
  string email = 3;
}
```

2. Generate code with:

`buf.gen.yaml`:

```yaml
version: v2
plugins:
  - local: [ "go",  "tool", "protoc-gen-go-fiber" ]
    out: out_dir
    opt:
      - paths=source_relative
      - error_handlers_package=github.com/petara94/protoc-gen-go-fiber/utils
      - json_unmarshal_package=encoding/json
      - grpc_error_handle_func=HandleGRPCStatusError
      - unmarshal_error_handle_func=HandleUnmarshalError
      - validation_error_handle_func=HandleValidationError
```

3. Register the routes in your Fiber app:

```go
app := fiber.New()
RegisterGreeterServiceFiberRoutes(app, serverImpl, grpcInterceptor)
```

## Supported Flags

| Flag                           | Description                                       | Default Value                                   |
|--------------------------------|---------------------------------------------------|-------------------------------------------------|
| `error_handlers_package`       | Path to the package with error handler functions  | `github.com/petara94/protoc-gen-go-fiber/utils` |
| `json_unmarshal_package`       | Path to the JSON unmarshal helper package         | `encoding/json`                                 |
| `parsers_package`              | Path to the package with parsers (for query/params)| `github.com/petara94/protoc-gen-go-fiber/utils` |
| `grpc_error_handle_func`       | Name of the gRPC error handler function           | `HandleGRPCStatusError`                         |
| `unmarshal_error_handle_func`  | Name of the JSON unmarshal error handler function | `HandleUnmarshalError`                          |
| `validation_error_handle_func` | Name of the validation error handler function     | `HandleValidationError`                         |

## Requirements

* Go 1.20+
* Required plugins:
    * `protoc-gen-validate`
    * `google/api/annotations.proto`
* Dependencies:
    * [Fiber](https://github.com/gofiber/fiber)
    * [gRPC](https://github.com/grpc/grpc-go)
    * [protoc-gen-validate](https://github.com/envoyproxy/protoc-gen-validate)

## Generated Code Structure

* Each method handler:
    * Parses the request body (if applicable).
    * Validates input (`Validate()`).
    * Injects HTTP headers into the context.
    * Calls the gRPC server method using the provided interceptor (if != nil).
    * Returns a JSON response or an error.

## License

[MIT License](LICENSE)

## Example: Using the POST /api/v1/user endpoint

You can now use the generated POST endpoint to create a user:

```bash
curl -X POST http://localhost:8080/api/v1/user \
  -H 'Content-Type: application/json' \
  -d '{"username": "john", "email": "john@example.com"}'
```

The response will be a JSON object with the created user's data.
