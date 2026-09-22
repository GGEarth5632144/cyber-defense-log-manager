# A Note Before We Begin

Thank you for taking the time to review my application and this assignment. I hope you enjoy reviewing my work, and I look forward to hearing from you soon.

# Full-Stack Log Management System

A multi-tenant log management system demo that ingests log data from multiple sources, normalizes it, and provides a web dashboard for visualization and alerting.

## Tech Stack
- **Backend**: Go (Fiber)
- **Database**: SQLite
- **Frontend**: React (Vite, TypeScript, TailwindCSS, Recharts)
- **Deployment**: Docker Compose

## Quick Start (Appliance Mode)

1. Ensure you have Docker and Docker Compose installed.
2. Run the system:
   ```bash
   make up
   # or
   docker-compose up -d
   ```
3. Access the dashboard at `http://localhost:3000`.
4. API is running at `http://localhost:8080`.
5. Syslog listener is running at `udp://localhost:514`.

## Documentation
- [Architecture](./docs/architecture.md)
- [Setup Appliance](./docs/setup_appliance.md)
- [Setup SaaS](./docs/setup_saas.md)
