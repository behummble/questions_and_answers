# Questions & Answers API Server

A simple web server for managing questions and answers with a RESTful API. This application allows users to create questions, add answers, and manage content through various API endpoints.

## Features

- **Questions Management**: Create, read, and delete questions
- **Answers Management**: Add, read, and delete answers for specific questions
- **RESTful API**: Clean HTTP endpoints for all operations
- **Docker Support**: Easy deployment using Docker Compose
- **Customizable Configuration**: Flexible configuration for server, database, and logging settings

## Prerequisites

- Docker
- Docker Compose

## Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/behummble/questions_and_answers.git
   cd questions_and_answers
   ```

2. **Install goose for migrations**
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```
   or for MacOs
   ```bash
   brew install goose 
   ```
3. **Setup configuration file (next topic)**
4. **Create .env file like example.env in root directory**

5. **Run with Docker Compose**
   ```bash
   docker-compose up -d
   ```
   or for MacOs
   ```bash
   docker compose up -d
   ```
6. **Up migrations**
   ```bash
   goose -dir ./migrations postgres "postgresql://your_username:your_password@your_host:your_port/your_dbName?sslmode=disable" up
   ```
   example:
   goose -dir ./migrations postgres "postgresql://admin:dGMavqa8-Z@127.0.0.1:5432/Questions?sslmode=disable" up


The API server will be available at the configured host and port (default: `localhost:8080`).

## Configuration

Change a `config.yaml` file in the project folder ./config to customize the application settings:

```yaml
server:
  host: "0.0.0.0"    # HTTP server host
  port: 8080         # HTTP server port

storage:
  host: "db"         # Database host (use "db" for Docker, "localhost" for local)
  port: 5432         # Database port
  name: "qa_db"      # Database name
  user: "postgres"   # Database username

log:
  level: 1      # Log level: debug, info, warning, error
  file: "app.log"    # Log file path
```

## Docker Compose

The `docker-compose.yml` file defines two services:

- **questions_answers**: The main API server
- **postgres**: PostgreSQL database
- **pg_admin**: PostgreSQL UI

### Environment Variables

You can override configuration using .env like in example.env:

## Logging

Logs are written to the configured log file and also output to stdout when running in Docker. Check the logs with:

```bash
docker-compose logs app
```

## Stopping the Application

```bash
docker-compose down
```

## API Documentation

The project includes comprehensive API documentation using Swagger/OpenAPI specification. Documetnation in docs/swagger.
