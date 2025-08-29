# Chat Agent

A Go-based chat agent application.

## Project Structure

```
.
├── main.go              # Main application entry point
├── cmd/                 # Command line applications
│   └── server/         # Server command
├── pkg/                # Public packages
│   ├── handlers/       # HTTP handlers
│   └── models/         # Data models
├── internal/           # Private packages
│   ├── config/         # Configuration management
│   └── database/       # Database operations
├── go.mod              # Go module definition
└── README.md           # This file
```

## Getting Started

### Prerequisites

- Go 1.23.4 or later

### Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```

### Running the Application

```bash
go run main.go
```

The server will start on `http://localhost:8080`

### Available Endpoints

- `GET /` - Welcome message
- `GET /health` - Health check endpoint

## Development

### Building

```bash
go build -o bin/chat-agent main.go
```

### Testing

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

