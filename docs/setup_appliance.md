# Appliance Setup Guide

Appliance mode is designed to run everything locally via Docker Compose.

## Prerequisites
- Docker
- Docker Compose

## Installation

1. **Clone the repository**:
   ```bash
   git clone <YOUR_GIT_REPO_URL>
   cd log-management-system
   ```
2. **Start the system**:
   ```bash
   docker-compose up -d --build
   ```
3. **Verify running containers**:
   ```bash
   docker-compose ps
   ```

## Usage

- **Frontend Dashboard**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Syslog Listener**: udp://localhost:514

## Testing

Run the included scripts in the `samples` folder to ingest data.

```bash
python3 samples/send_syslog.py
python3 samples/post_logs.py
```
