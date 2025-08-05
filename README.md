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
   docker run --name pharmacy-db -e POSTGRES_PASSWORD=your_password -e POSTGRES_DB=pharmacy_claims -p 5432:5432 -d postgres
   ```
   
   The PostgreSQL Docker image can be found on Docker Hub: [https://hub.docker.com/_/postgres](https://hub.docker.com/_/postgres)

4. Run the application:
   ```bash
   go run main.go
   ```

## Project Structure

- `main.go` - Main application entry point

## Requirements

- Go 1.19 or higher
- Docker (latest version)
- golang-migrate (latest version)
- sqlc (latest version) 