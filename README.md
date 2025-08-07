# Pharmacy Claims Application

This is a Go application for managing pharmacy claims.

## Prerequisites

This project requires the following tools to be installed on your system:

- **Go** - Programming language
- **Docker** - Containerization platform
- **golang-migrate** - Database migration tool
- **sqlc** - SQL compiler and code generator

### Installing Go

1. **Download Go**: Visit the official Go download page at [https://golang.org/dl/](https://golang.org/dl/)

2. **Choose your platform**:
   - **macOS**: Download the `.pkg` installer for macOS
   - **Windows**: Download the `.msi` installer for Windows
   - **Linux**: Download the appropriate `.tar.gz` file for your Linux distribution

3. **Install Go**:
   - **macOS**: Run the downloaded `.pkg` file and follow the installation wizard
   - **Windows**: Run the downloaded `.msi` file and follow the installation wizard
   - **Linux**: Extract the `.tar.gz` file to `/usr/local` and add Go to your PATH

4. **Verify installation**: Open a terminal/command prompt and run:
   ```bash
   go version
   ```

   You should see output similar to:
   ```
   go version go1.21.0 darwin/amd64
   ```

### Installing Docker

1. **Download Docker**: Visit the official Docker download page at [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)

2. **Choose your platform**:
   - **macOS**: Download Docker Desktop for Mac
   - **Windows**: Download Docker Desktop for Windows
   - **Linux**: Follow the installation guide for your specific distribution

3. **Install Docker**:
   - **macOS**: Run the downloaded `.dmg` file and drag Docker to Applications
   - **Windows**: Run the downloaded `.exe` file and follow the installation wizard
   - **Linux**: Follow the distribution-specific installation commands

4. **Verify installation**: Open a terminal/command prompt and run:
   ```bash
   docker --version
   ```

   You should see output similar to:
   ```
   Docker version 24.0.0, build af53cfe
   ```

### Installing golang-migrate

1. **Using Go install** (recommended):
   ```bash
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ```

2. **Using package managers**:
   - **macOS** (using Homebrew):
     ```bash
     brew install golang-migrate
     ```
   - **Linux** (using apt):
     ```bash
     sudo apt-get install golang-migrate
     ```

3. **Verify installation**: Open a terminal/command prompt and run:
   ```bash
   migrate -version
   ```

   You should see output similar to:
   ```
   migrate version v4.16.2
   ```

### Installing sqlc

1. **Using Go install** (recommended):
   ```bash
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ```

2. **Using package managers**:
   - **macOS** (using Homebrew):
     ```bash
     brew install sqlc
     ```
   - **Ubuntu** (using snap):
     ```bash
     sudo snap install sqlc
     ```

3. **Using Docker**:
   ```bash
   docker pull sqlc/sqlc
   ```

4. **Verify installation**: Open a terminal/command prompt and run:
   ```bash
   sqlc version
   ```

   You should see output similar to:
   ```
   sqlc version v1.29.0
   ```

   For more installation options, visit the official sqlc documentation: [https://docs.sqlc.dev/en/stable/overview/install.html](https://docs.sqlc.dev/en/stable/overview/install.html)

## Getting Started

1. Clone this repository
2. Navigate to the project directory
3. Start the PostgreSQL database using Docker:
   ```bash
   docker run --name pharmacy-db -e POSTGRES_PASSWORD=your_password -e POSTGRES_USER=root -e POSTGRES_DB=pharmacy_claims -p 5432:5432 -d postgres
   ```
   
   The PostgreSQL Docker image can be found on Docker Hub: [https://hub.docker.com/_/postgres](https://hub.docker.com/_/postgres)


4. Run database migrations:
    ```migrate -path db/migration -database "postgresql://root:your_password@localhost:5432/pharmacy_claims?sslmode=disable" -verbose up
     ```

5. Generate SQL code:
   ```bash
   sqlc generate
   ```

8. Run the application:
   ```bash
   go run main.go
   ```

## API Documentation

The application provides a RESTful API for managing pharmacy claims.

### Base URL
```
http://localhost:8080
```

### Endpoints

#### Health Check
- **GET** `/health`
- Returns server health status

#### Claims Management

**Create Claim**
- **POST** `/api/v1/claims`
- **Body:**
  ```json
  {
    "ndc": "123456789",
    "quantity": 30,
    "npi": "9876543210",
    "price": 15.99
  }
  ```
- **Response:**
  ```json
  {
    "status": "claim submitted",
    "claim_id": "abc123"
  }
  ```

**Get Claim**
- **GET** `/api/v1/claims/{id}`
- **Response:**
  ```json
  {
    "status": "claim retrieved",
    "claim_id": "abc123",
    "claim": {
      "id": "abc123",
      "ndc": "123456789",
      "quantity": 30,
      "npi": "9876543210",
      "price": 15.99,
      "timestamp": "2024-01-01T12:00:00Z"
    }
  }
  ```

**Create Reversal**
- **POST** `/api/v1/reversals`
- **Body:**
  ```json
  {
    "claim_id": "abc123"
  }
  ```
- **Response:**
  ```json
  {
    "status": "claim reversed",
    "claim_id": "abc123"
  }
  ```

## Event Logging

All claim submissions and reversals are automatically logged to `logs/pharmacy_events.json` in JSON format.

### Event Log Format

**Claim Submission Event:**
```json
{
  "id": "event-uuid",
  "type": "claim_submitted",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "claim_id": "claim-uuid",
    "ndc": "123456789",
    "npi": "9876543210",
    "quantity": 30,
    "price": 15.99
  }
}
```

**Claim Reversal Event:**
```json
{
  "id": "event-uuid",
  "type": "claim_reversed",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "claim_id": "claim-uuid"
  }
}
```  

**Get Claim**
- **GET** `/api/v1/claims/{id}`
- Returns a specific claim by ID

**Create Reversal**
- **POST** `/api/v1/reversals`
- **Body:** 
  ```json
  {
   "claim_id": "abc123"
  }


### Response Format
```json
{
  "status": "claim reversed",
  "claim_id": "abc123"
}
  ```

### Error Response Format
```json
{
  "status": "error",
  "message": "Specific error description",
  "code": 400
}
```

### Common Error Examples

**Invalid JSON:**
```json
{
  "status": "error",
  "message": "Invalid JSON format in request body",
  "code": 400
}
```

**Missing Required Fields:**
```json
{
  "status": "error",
  "message": "NDC (National Drug Code) is required",
  "code": 400
}
```

**Invalid UUID:**
```json
{
  "status": "error",
  "message": "Invalid claim ID format. Must be a valid UUID",
  "code": 400
}
```

**Claim Not Found:**
```json
{
  "status": "error",
  "message": "Claim not found",
  "code": 404
}
```

## Project Structure

- `main.go` - Main application entry point
- `env.example` - Environment variables template

## Environment Variables

This project uses environment variables for database configuration. You can set them in several ways:

1. **Create a `.env` file** (recommended):
   ```bash
   cp env.example .env
   # Edit .env file with your credentials
   ```

2. **Export environment variables**:
   ```bash
   export DB_PASSWORD=your_secure_password
   export DB_USER=root
   export DB_NAME=pharmacy_claims
   ```

**Security Note**: Never commit passwords to version control. The `.env` file is already in `.gitignore`.

## Requirements

- Go 1.19 or higher
- Docker (latest version)
- golang-migrate (latest version)
- sqlc (latest version) 