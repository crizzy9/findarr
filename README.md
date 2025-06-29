# Findarr

A universal content metadata resolver, content finder, and content organizer for self-hosted systems. Findarr can handle various media types including movies, shows, books, audiobooks, music, and more.

## Features

- Modern architecture with Go backend and htmx frontend
- Configurable metadata sources for different types of media
- Clean, responsive UI
- Easy navigation and organization of your media collection
- Inspired by *arr applications (Radarr, Sonarr, Readarr) but with a universal approach

## Project Structure

```
findarr/
├── cmd/
│   └── server/         # main.go entry point
├── internal/
│   ├── api/            # HTTP handlers/controllers
│   ├── metadata/       # Metadata fetcher/resolver interfaces & implementations
│   ├── storage/        # Storage (DB, filesystem, etc.)
│   └── config/         # Config management
├── web/
│   ├── static/         # JS, CSS, images
│   └── templates/      # HTML templates (htmx-ready)
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

### Prerequisites

- Option 1 (Docker):
  - Docker and Docker Compose

- Option 2 (Local Development):
  - Go 1.20 or higher
  - Git

### Installation

#### Option 1: Docker (Recommended)

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/findarr.git
   cd findarr
   ```

2. Use the setup script (easiest):
   ```
   ./setup.sh start
   ```
   
   Or use Make:
   ```
   make docker-run
   ```
   
   Or manually with Docker Compose:
   ```
   docker-compose up -d
   ```

3. Open your browser and navigate to `http://localhost:8080`

4. To stop the container:
   ```
   ./setup.sh stop
   ```
   
   Or use Make:
   ```
   make docker-stop
   ```
   
   Or manually:
   ```
   docker-compose down
   ```

5. Other setup script commands:
   ```
   ./setup.sh logs     # Show container logs
   ./setup.sh restart  # Restart the container
   ./setup.sh status   # Show container status
   ./setup.sh build    # Build the Docker image
   ./setup.sh help     # Show help message
   ```

#### Option 2: Local Development

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/findarr.git
   cd findarr
   ```

2. Install dependencies:
   ```
   go get -u github.com/gorilla/mux
   ```

3. Run the application:
   ```
   make run
   ```
   Or manually:
   ```
   go run cmd/server/main.go
   ```

4. Open your browser and navigate to `http://localhost:8080`

### Configuration

Findarr can be configured using environment variables:

| Environment Variable | Description | Default |
|----------------------|-------------|---------|
| FINDARR_HOST | Host to bind the server to | 0.0.0.0 |
| FINDARR_PORT | Port to bind the server to | 8080 |
| FINDARR_DB_PATH | Path to the database file | config/findarr.db |
| FINDARR_MOVIES_PATH | Path to movies directory | "" |
| FINDARR_SHOWS_PATH | Path to TV shows directory | "" |
| FINDARR_BOOKS_PATH | Path to books directory | "" |
| FINDARR_MUSIC_PATH | Path to music directory | "" |

When using Docker, you can set these in the `docker-compose.yml` file.

## Development

### Adding a New Metadata Provider

1. Implement the `metadata.Provider` interface in a new file under `internal/metadata/`
2. Register the provider in the application

### Adding a New Media Type

1. Define the media type structure in `internal/metadata/types.go`
2. Create appropriate handlers and templates

## Docker Compose Configuration

The included Docker Compose setup provides an easy way to self-host Findarr. Here's how to customize it:

### Environment Variables

Edit the `docker-compose.yml` file to set environment variables:

```yaml
environment:
  - TZ=UTC  # Set your timezone
  - FINDARR_PORT=8080  # Server port
  - FINDARR_HOST=0.0.0.0  # Server host
  - FINDARR_DB_PATH=/app/config/findarr.db  # Database path
  - FINDARR_MOVIES_PATH=/media/movies  # Path to movies
  - FINDARR_SHOWS_PATH=/media/shows  # Path to TV shows
  - FINDARR_BOOKS_PATH=/media/books  # Path to books
  - FINDARR_MUSIC_PATH=/media/music  # Path to music
```

### Volume Mounts

Mount your media directories and configuration:

```yaml
volumes:
  - ./config:/app/config  # For configuration persistence
  - /path/to/movies:/media/movies  # Your movies directory
  - /path/to/shows:/media/shows  # Your TV shows directory
  - /path/to/books:/media/books  # Your books directory
  - /path/to/music:/media/music  # Your music directory
```

### Port Mapping

Change the port mapping if needed:

```yaml
ports:
  - "8080:8080"  # host:container
```

## License

MIT

## Acknowledgments

- Inspired by Radarr, Sonarr, Readarr, and other *arr applications
- Built with Go, htmx, and Tailwind CSS
