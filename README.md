# rivertui

A terminal-based user interface and CLI for [River Queue](https://riverqueue.com/) - monitor and manage your jobs from the command line.

<div align="center">
  <img src="./t-rec.gif" alt="t-rec">
</div>

## Installation

```bash
go install github.com/almottier/rivertui@latest
```

## Usage

```bash
export RIVER_DATABASE_URL="postgres://localhost:5432/myapp"
rivertui
```

### Command Line Options

| Flag | Environment Variable | Description | Default |
|------|---------------------|-------------|---------|
| `--database-url` | `RIVER_DATABASE_URL` | PostgreSQL connection string | Required |
| `--refresh` | - | Refresh interval | `1s` |
| `--job-id` | - | Start in details view for specific job ID | - |

### Example

```bash
# Basic usage
rivertui --database-url "postgres://localhost:5432/myapp"

# Custom refresh rate
rivertui --database-url "postgres://localhost:5432/myapp" --refresh 2s

# Start viewing specific job
rivertui --database-url "postgres://localhost:5432/myapp" --job-id 12345

# Using environment variable
export RIVER_DATABASE_URL="postgres://localhost:5432/myapp"
rivertui
```

## Features

- **Real-time job monitoring** with auto-refresh
- **Job state filtering** (available, running, completed, discarded, etc.)
- **Job kind filtering** and search
- **Job details view** with full arguments, metadata, and error information
- **Job operations**: retry and cancel jobs
- **Pagination** for large job lists
- **Queue management**: TODO
- **Keyboard-driven navigation**

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Enter` | View job details |
| `/` | Search by job kind or jump to job ID |
| `0-7` | Filter by job state (0=All, 1=Completed, 2=Available, etc.) |
| `r` | Retry selected job |
| `c` | Cancel selected job |
| `n` | Next page |
| `p` | Previous page |
| `q` | Quit |

## Requirements

- Go 1.21+
- River Queue 0.5.0+
- PostgreSQL database with River Queue tables
- Terminal with color support

## License

MIT
