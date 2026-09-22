# Log Manager Backend API

This directory contains the core Go application that serves as both the REST API for the frontend dashboard and the data ingestion processor. 

## Architectural Design
We chose **Go (Golang)** with the **Fiber** framework for the backend to ensure a lightweight, high-performance footprint that can easily run on a single cloud VM (Appliance mode). 

For storage, we utilize **SQLite**. SQLite was chosen because it requires zero external dependencies, making the Docker deployment incredibly simple. However, because SQLite does not handle concurrent multi-process writes well, we deliberately bundled the **Ingestion Engine** and the **API Engine** into this single monolithic Go service. This ensures a single database connection pool manages all read/write locks, preventing database corruption or locking crashes.

## Core Components
- **`main.go`**: Bootstraps the HTTP server and the background workers. It also contains the raw UDP Syslog listener (Port 514).
- **`api.go`**: Contains all authenticated REST endpoints (`/api/stats`, `/api/logs`, etc.) used by the React frontend. It handles dynamic SQL generation for tenant and time-range filters.
- **`database.go`**: Manages the SQLite connection pool and the automated table schemas.
- **`alerting.go`**: A background goroutine that polls the database every 60 seconds to detect threshold violations (e.g., repeated failed logins).
- **`retention.go`**: A background goroutine that executes a cleanup query every hour, enforcing the 7-day data retention policy.
- **`models.go`**: Defines the shared schema and struct mappings.

## Authentication & Multi-Tenancy
The API uses **JWT (JSON Web Tokens)** for stateless authentication. When a user logs in, the `AuthMiddleware` intercepts subsequent requests, extracts the `tenant` claim from the token, and securely appends `AND tenant = ?` to all SQL queries, guaranteeing strict data isolation between tenants.
