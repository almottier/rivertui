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

| Flag             | Environment Variable | Description                               | Default  |
| ---------------- | -------------------- | ----------------------------------------- | -------- |
| `--database-url` | `RIVER_DATABASE_URL` | PostgreSQL connection string              | Required |
| `--refresh`      | -                    | Refresh interval                          | `1s`     |
| `--job-id`       | -                    | Start in details view for specific job ID | -        |
| `--kind`         | -                    | Start with kind filter applied            | -        |

### Example

```bash
# Basic usage
rivertui --database-url "postgres://localhost:5432/myapp"

# Custom refresh rate
rivertui --database-url "postgres://localhost:5432/myapp" --refresh 2s

# Start viewing specific job
rivertui --database-url "postgres://localhost:5432/myapp" --job-id 12345

# Start with kind filter applied
rivertui --database-url "postgres://localhost:5432/myapp" --kind "SendEmailJob"

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
- **Queue management**: view, pause, and resume queues
- **Keyboard-driven navigation**

## Keyboard Shortcuts

| Key      | Action                                                      |
| -------- | ----------------------------------------------------------- |
| `Enter`  | View job details                                            |
| `/`      | Search by job kind or jump to job ID                        |
| `0-7`    | Filter by job state (0=All, 1=Completed, 2=Available, etc.) |
| `Ctrl+Q` | View queues                                                 |
| `r`      | Retry selected job                                          |
| `c`      | Cancel selected job                                         |
| `n`      | Next page                                                   |
| `p`      | Previous page                                               |
| `p`      | Pause selected queue                                        |
| `r`      | Resume selected queue                                       |
| `q`      | Quit                                                        |

## Color Themes & Customization

Custom color schemes can be set via the following `RIVER_COLOR_` prefixed environment variables using hex color strings:

```bash
# Example: Gruvbox Dark
export RIVER_COLOR_PRIMARY="#fbf1c7"                  # Gruvbox Dark fg0
export RIVER_COLOR_SECONDARY="#ebdbb2"                # Gruvbox Dark fg1
export RIVER_COLOR_TERTIARY="#d5c4a1"                 # Gruvbox Dark fg2
export RIVER_COLOR_BORDER="#665c54"                   # Gruvbox Dark bg3
export RIVER_COLOR_WARNING="#d79921"                  # Yellow
export RIVER_COLOR_INFO="#458588"                     # Blue
export RIVER_COLOR_SUCCESS="#98971a"                  # Green
export RIVER_COLOR_ERROR="#cc241d"                    # Red
export RIVER_COLOR_AVAILABLE="#689d6a"                # Aqua
export RIVER_COLOR_CANCELLED="#d65d0e"                # Orange
export RIVER_COLOR_RETRYABLE="#b16286"                # Purple
export RIVER_COLOR_SCHEDULED="#83a598"                # Muted Blue variant
export RIVER_COLOR_TITLE="#fb4934"                    # Bright Red (for headings)
export RIVER_COLOR_CONTRAST_SECONDARY="#3c3836"       # Gruvbox Dark bg1
export RIVER_COLOR_SELECTED_BG="#504945"              # Gruvbox Dark bg2
export RIVER_COLOR_CONTRAST_BACKGROUND="#282828"      # Gruvbox Dark bg0
export RIVER_COLOR_PRIMATIVE_BACKGROUND="#1d2021"     # Even darker for depth
export RIVER_COLOR_MORE_CONTRAST_BACKGROUND="#141617" # Extra dark contrast
export RIVER_COLOR_SELECTED_FG="#fbf1c7"              # Gruvbox Dark text

# you can also use a boolean flag to make the background transparent
export RIVER_COLOR_TRANSPARENT_BG=true
```

## Requirements

- Go 1.21+
- River Queue 0.5.0+
- PostgreSQL database with River Queue tables
- Terminal with color support

## License

MIT
