# Discovery Tree

This repository contains a discovery tree - a structured way to organize and explore knowledge, ideas, or information in a hierarchical format. It serves as a visual map for understanding relationships between different concepts and can be used for research, learning, or project planning purposes.

## REST API

The Discovery Tree includes a REST API for managing tasks in a hierarchical tree structure.

### Running the API

To start the API server:

```bash
# Run directly
go run cmd/api/main.go

# Or build and run
go build -o discovery-tree-api cmd/api/main.go
./discovery-tree-api
```

### Configuration

The API can be configured using environment variables. All configuration options have sensible defaults for development.

| Environment Variable | Default Value | Description |
|---------------------|---------------|-------------|
| `PORT` | `8080` | Port number for the HTTP server |
| `DATA_PATH` | `./data/tasks.json` | Path to the JSON file for task persistence |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |
| `ENABLE_CORS` | `true` | Enable Cross-Origin Resource Sharing |
| `ENABLE_SWAGGER` | `true` | Enable Swagger/OpenAPI documentation |

### Example Configuration

```bash
# Set custom port and data path
export PORT=3000
export DATA_PATH=/var/data/tasks.json
export LOG_LEVEL=debug

# Start the server
go run cmd/api/main.go
```

### API Documentation

When `ENABLE_SWAGGER` is true (default), interactive API documentation is available at:
- Swagger UI: `http://localhost:8080/api/v1/docs`
- OpenAPI Schema: `http://localhost:8080/api/v1/swagger.json`

### Health Check

The API provides a health check endpoint at:
- `GET /health` - Returns server status
