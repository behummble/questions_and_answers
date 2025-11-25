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
   git clone <repository-url>
   cd <project-directory>
   ```

2. **Run with Docker Compose**
   ```bash
   docker-compose up -d
   ```

3. **Verify the server is running**
   ```bash
   curl http://localhost:8080/health
   ```

The API server will be available at the configured host and port (default: `localhost:8080`).

## Configuration

Create a `config.yaml` file in the project root to customize the application settings:

```yaml
server:
  host: "0.0.0.0"    # HTTP server host
  port: 8080         # HTTP server port

database:
  host: "db"         # Database host (use "db" for Docker, "localhost" for local)
  port: 5432         # Database port
  name: "qa_db"      # Database name
  user: "postgres"   # Database username
  password: "password" # Database password

logging:
  level: "info"      # Log level: debug, info, warning, error
  file: "app.log"    # Log file path
```

## API Endpoints

### Questions

- **Create Question**
  - `POST /api/questions`
  - Body: `{"question": "Your question text"}`

- **Get All Questions**
  - `GET /api/questions`

- **Get Specific Question**
  - `GET /api/questions/{question_id}`

- **Delete Question**
  - `DELETE /api/questions/{question_id}`

### Answers

- **Create Answer**
  - `POST /api/questions/{question_id}/answers`
  - Body: `{"answer": "Your answer text"}`

- **Get Specific Answer**
  - `GET /api/answers/{answer_id}`

- **Delete Answer**
  - `DELETE /api/answers/{answer_id}`

## Example Usage

### Create a Question
```bash
curl -X POST http://localhost:8080/api/questions \
  -H "Content-Type: application/json" \
  -d '{"question": "What is Docker?"}'
```

### Add an Answer to a Question
```bash
curl -X POST http://localhost:8080/api/questions/1/answers \
  -H "Content-Type: application/json" \
  -d '{"answer": "Docker is a containerization platform."}'
```

### Get All Questions with Answers
```bash
curl http://localhost:8080/api/questions
```

## Docker Compose

The `docker-compose.yml` file defines two services:

- **app**: The main API server
- **db**: PostgreSQL database

### Environment Variables

You can override configuration using environment variables:

```yaml
environment:
  - SERVER_HOST=0.0.0.0
  - SERVER_PORT=8080
  - DB_HOST=db
  - DB_PORT=5432
  - DB_NAME=qa_db
  - DB_USER=postgres
  - DB_PASSWORD=password
  - LOG_LEVEL=info
```

## Development

### Running without Docker

1. Install dependencies
2. Ensure PostgreSQL is running
3. Update `config.yaml` with local database credentials
4. Run the application:
   ```bash
   ./start-server.sh
   ```

### Building Docker Image Manually

```bash
docker build -t qa-server .
docker run -p 8080:8080 qa-server
```

## Logging

Logs are written to the configured log file and also output to stdout when running in Docker. Check the logs with:

```bash
docker-compose logs app
```

## Stopping the Application

```bash
docker-compose down
```

## License

[Add your license information here]