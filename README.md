# cfgctl

A lightweight CLI tool for creating and managing simple configuration files.

## Features

- Initialize a configuration file
- Set or update key/value pairs
- Retrieve values by key
- List all stored configuration entries
- Delete configuration keys
- Includes mock storage and test utilities for unit testing

---

## Project Structure

```text
cfgctl/
├── cmd/
│   └── cfgctl/
│       └── main.go              # CLI entrypoint
│
├── internal/
│   ├── app/
│   │   └── cfgctl/
│   │       └── commands/        # Command implementations
│   │           ├── init.go
│   │           ├── set.go
│   │           ├── get.go
│   │           ├── list.go
│   │           ├── delete.go
│   │           ├── errors.go
│   │           └── *_test.go
│   │
│   ├── pkg/
│   │   └── storage/             # Storage interface and mock storage
│   │
│   └── testutil/                # Shared test helpers
│
├── bin/                         # Compiled binaries
├── testbin/                     # Test binaries
├── tmp/                         # Temporary files
├── Makefile
├── go.mod
└── README.md
```

---

## Requirements

- Go 1.20 or newer
- Git
- Make (optional)

---

## Installation

### Clone the Repository

```bash
git clone https://github.com/FeelsCoderMan/cfgctl.git
cd cfgctl
```

### Download Dependencies

```bash
go mod download
```

### Build the Project

Using Make:

```bash
make build
```

Or directly with Go:

```bash
go build -o bin/cfgctl cmd/cfgctl/main.go
```

---

## Usage

### Initialize a Configuration File

Creates a new configuration file at the default path.

```bash
cfgctl init
```

### Set a Key

```bash
cfgctl set <key> <value>
```

Example:

```bash
cfgctl set username admin
```

### Get a Key

```bash
cfgctl get <key>
```

Example:

```bash
cfgctl get username
```

### List All Keys

```bash
cfgctl list
```

### Delete a Key

```bash
cfgctl delete <key>
```

Example:

```bash
cfgctl delete username
```

### Help

It is not implemented yet.
```bash
cfgctl --help
```

---

## Testing

Run all unit tests using Make:

```bash
make test
```

Or directly with Go:

```bash
go test ./...
```

The project includes:

- Mock storage implementations
- File-based test helpers
- Deterministic command testing utilities

---

## Makefile Targets

| Command        | Description                     |
|----------------|---------------------------------|
| `make build`   | Build the CLI binary            |
| `make test`    | Run all tests                   |
| `make clean`   | Remove build artifacts          |

---

## Output Directories

| Directory  | Purpose                    |
|------------|----------------------------|
| `bin/`     | Compiled application files |
| `testbin/` | Test binaries              |
| `tmp/`     | Temporary runtime files    |

---

## Example Workflow

```bash
# Initialize config
cfgctl init

# Set values
cfgctl set host localhost
cfgctl set port 8080

# Retrieve a value
cfgctl get host

# List all entries
cfgctl list

# Delete a value
cfgctl delete port
```

---

## License

This project is open source and available under the MIT License.
