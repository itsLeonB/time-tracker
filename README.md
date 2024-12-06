# Time Tracker

An API for tracking time spent working on tasks in a project.


## Getting Started

### Prerequisites

- Go 1.16+
- Docker (for containerization)
- PostgreSQL

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/itsLeonB/time-tracker.git
    cd time-tracker
    ```

2. Install dependencies:
    ```sh
    go mod download
    ```

3. Set up environment variables:
    ```sh
    cp .env.example .env
    ```

4. Run database migrations in `./db/migrations/up.sql`

### Running the Application

To start the application, run:
```sh
go run cmd/app/main.go
```
The application will be available at `http://localhost:8000`

Or simply run with docker:
```sh
docker compose up
```

## API Documentation

The documentation can be viewed [here](https://documenter.getpostman.com/view/32713619/2sAYBbeUzT).