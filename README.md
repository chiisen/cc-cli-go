# Claude Code CLI - Go Implementation

A Go implementation of Claude Code CLI for learning purposes. This project implements a minimal but functional CLI tool that interacts with Anthropic's Claude API.

## Project Location

```
/Users/liao-eli/github/cc-cli/cc-cli-go/
```

## Features

### Core Features
- ✅ CLI Entry & Argument Parsing (using Cobra)
- ✅ Anthropic API Streaming Client (SSE)
- ✅ Query Loop (Channel-based streaming)
- ✅ Tool System (Interface-based, extensible)
- ✅ Bash Tool (Execute shell commands)
- ✅ Read Tool (Read file contents)
- ✅ Edit Tool (Edit files by exact string replacement)
- ✅ Basic TUI (Bubble Tea framework)
- ✅ Message Rendering
- ✅ Streaming Event Handling
- ✅ Concurrent Tool Execution

### Technical Highlights
- **Idiomatic Go**: Follows Go best practices and conventions
- **Channel-based Streaming**: Uses Go channels instead of generators
- **Interface-based Design**: Tool system is easily extensible
- **Concurrent Execution**: Parallel tool execution with goroutines + WaitGroup
- **Single Binary**: No runtime dependencies

## Directory Structure

```
cc-cli-go/
├── cmd/
│   └── claude-code/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/                     # Anthropic API client
│   │   ├── client.go            # HTTP client setup
│   │   ├── config.go            # API configuration
│   │   ├── events.go            # Event parsing
│   │   ├── request.go           # Request builder
│   │   └── streaming.go         # SSE streaming
│   ├── cli/                     # CLI framework
│   │   ├── root.go              # Root command
│   │   └── run.go               # Interactive mode command
│   ├── query/                   # Query engine
│   │   ├── engine.go            # Core query loop
│   │   └── types.go             # Query types
│   ├── tools/                   # Tool system
│   │   ├── tool.go              # Tool interface
│   │   ├── registry.go          # Tool registry
│   │   ├── bash/
│   │   │   └── bash.go          # Bash tool
│   │   ├── edit/
│   │   │   └── edit.go          # Edit tool
│   │   └── read/
│   │       └── read.go          # Read tool
│   ├── tui/                     # Terminal UI
│   │   ├── app.go               # Bubble Tea app
│   │   └── model.go             # TUI model
│   └── types/                   # Shared types
│       ├── content.go           # Content block types
│       ├── message.go           # Message types
│       ├── usage.go             # Token usage
│       └── uuid.go              # UUID generation
├── bin/                         # Compiled binaries
│   └── cc-cli-go
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
└── .gitignore
```

## Prerequisites

- Go 1.21 or higher
- Anthropic API Key

## Installation

### Clone the Repository

```bash
cd /Users/liao-eli/github/cc-cli
git clone <repository-url> cc-cli-go
cd cc-cli-go
```

### Install Dependencies

```bash
go mod download
```

### Build

```bash
# Build binary
go build -o bin/cc-cli-go ./cmd/claude-code

# Or use go install
go install ./cmd/claude-code
```

## Usage

### Set API Key

```bash
export ANTHROPIC_API_KEY="your-api-key-here"
```

### Run Interactive Mode

```bash
# Using compiled binary
./bin/cc-cli-go run

# Or using go run
go run ./cmd/claude-code run
```

### Check Version

```bash
./bin/cc-cli-go --version
# Output: claude-code version 0.1.0
```

### Available Commands

```bash
./bin/cc-cli-go --help
```

## Testing

### Build Test

```bash
# Build all packages
go build ./...

# Build specific package
go build ./internal/api
go build ./internal/tools
go build ./internal/query
```

### Run Tests

```bash
# Run all tests (when available)
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/tools/bash
```

### Manual Testing

1. **Test Version Command**
   ```bash
   ./bin/cc-cli-go --version
   ```

2. **Test Interactive Mode**
   ```bash
   ANTHROPIC_API_KEY=your-key ./bin/cc-cli-go run
   ```
   
   - Type a message and press Enter
   - Watch the streaming response
   - Press Ctrl+C to exit

3. **Test Tools (when API is available)**
   - The assistant can use Bash, Read, and Edit tools
   - Tools are executed concurrently when safe
   - Results are displayed in the TUI

## Development

### Project Statistics

- **Total Files**: 22 Go source files
- **Total Lines**: ~1,500+ LOC
- **Packages**: 8 internal packages
- **Dependencies**: 
  - github.com/spf13/cobra (CLI)
  - github.com/charmbracelet/bubbletea (TUI)
  - github.com/charmbracelet/bubbles (TUI components)
  - github.com/charmbracelet/lipgloss (styling)
  - github.com/google/uuid (UUID generation)

### Architecture

```
User Input → TUI (Bubble Tea) → Query Engine → API Client → Claude API
                ↓                      ↓
            Messages            Tool Execution
                                       ↓
                              Tool Registry
                                       ↓
                            Bash / Read / Edit
```

### Key Design Decisions

1. **Channel-based Streaming**: Uses Go channels for event streaming instead of TypeScript generators
2. **Interface-based Tools**: Tool interface allows easy addition of new tools
3. **Concurrent Execution**: Tools are executed in parallel when concurrency-safe
4. **Explicit Error Handling**: Go's explicit error handling pattern

## Future Enhancements

Potential features for future development:

- [ ] Permission System (tool approval dialogs)
- [ ] Glob Tool (file pattern matching)
- [ ] Grep Tool (search file contents)
- [ ] Write Tool (create new files)
- [ ] Session Storage (conversation persistence)
- [ ] CLAUDE.md Discovery (automatic context loading)
- [ ] Configuration Management (settings files)
- [ ] Improved Error Handling
- [ ] Comprehensive Test Suite
- [ ] Documentation Generation

## Troubleshooting

### Build Errors

```bash
# Clean build cache
go clean -cache

# Update dependencies
go mod tidy
go mod download

# Rebuild
go build -o bin/cc-cli-go ./cmd/claude-code
```

### API Key Issues

```bash
# Verify API key is set
echo $ANTHROPIC_API_KEY

# Set API key
export ANTHROPIC_API_KEY="sk-ant-api03-..."
```

### Runtime Errors

- **"ANTHROPIC_API_KEY environment variable is required"**: Set your API key
- **Connection errors**: Check network connectivity
- **Rate limiting**: Wait and retry with exponential backoff

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request

## License

This project is for learning purposes. See LICENSE file for details.

## Acknowledgments

- Original Claude Code CLI by Anthropic
- Bubble Tea framework by Charm
- Cobra CLI framework by Spf13

## Version History

- **v0.1.0** (2026-04-01): Initial release with core features
  - Basic CLI structure
  - API streaming client
  - Tool system (Bash, Read, Edit)
  - Query engine
  - Basic TUI

---

**Note**: This is a learning project and not intended for production use. For production use, please refer to the official Claude Code CLI.