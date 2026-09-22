# LogManager

A demo multi-tenant log management system. It collects security logs over syslog and HTTPS, normalizes them to one schema, and lets each tenant search them, chart them, and receive alerts. The same Compose stack runs as an appliance on one machine and as a SaaS deployment on a cloud VM.

**Ingestion:** Syslog over UDP on port 514, and JSON logs over HTTP via `POST /ingest` on port 8080.  
**Normalization:** Heuristic syslog parsing, and JSON mapping for AWS CloudTrail, Microsoft 365, Active Directory, CrowdStrike, and custom apps. All are mapped to a single unified event schema, keeping the original payload intact.  
**Storage and search:** SQLite optimized with indexes, featuring an automated background engine enforcing 7-day retention.  
**Dashboard:** A visual timeline and aggregated metrics for top source IPs, users, and event types, dynamically filtered by time range and tenant.  
**Alerting:** Background threshold rules evaluated every minute (e.g., repeated failed logins from the same IP). Alerts are immediately surfaced in the UI.  
**Security:** JWT-based stateless sign-in, strict `admin` and `viewer` roles, and robust tenant data isolation enforced directly at the API middleware layer.  

## Live Demo
The SaaS deployment is accessible at `http://<your-cloud-vm-ip>`. 
A walkthrough of the architecture, ingestion, search, dashboard, and alerting is provided in the official 30-minute demo video.

## Quick Start
You need Linux, macOS, or Windows with Docker Engine and the Compose plugin, plus `python3` for sending samples. Ports 3000, 8080, and 514 must be free.

```bash
git clone https://github.com/your-username/log-management-system.git
cd log-management-system
make up     # build and start the stack in the background
```

Open `http://localhost:3000`. The system automatically seeds the following test accounts:

| Account Username | Role | Sees | Password |
| :--- | :--- | :--- | :--- |
| `admin` | admin | every tenant (`*`) | `admin123` |
| `viewer` | viewer | strictly `demoA` | `viewer123` |
| `viewer_b` | viewer | strictly `demoB` | `viewer123` |

Once the stack is running, simulate live traffic using the provided scripts:
```bash
make test                      # Run integration tests to verify API health
python samples/send_syslog.py  # Fire simulated UDP syslog traffic
python samples/post_logs.py    # Fire simulated AWS/M365/CrowdStrike HTTP traffic
```

## Common Tasks
```bash
make up          # Start the stack and rebuild containers if code changed
make down        # Stop the stack and securely delete the database volume
make logs        # Tail the logs of all running services
make test        # Run the Python API integration tests against the live stack
```

## Project Structure

| Path | What |
| :--- | :--- |
| `docker-compose.yml`, `Makefile`, `.env.example` | The appliance stack definition, common task runner, and environment variable templates. |
| [`backend/`](backend/) | The monolithic Go service. Handles the REST API (Fiber), background alerting, retention, and direct SQLite interactions. |
| [`frontend/`](frontend/) | The React 18 SPA. Built with Vite, Tailwind CSS, and Recharts. Served by Nginx. |
| [`ingest/`](ingest/) | Documentation outlining the ingestion architecture. (Note: To ensure SQLite stability, the actual ingestion logic runs tightly coupled within `/backend/ingest.go` and `/backend/main.go`). |
| [`samples/`](samples/) | Sample JSON logs, syslog examples, and the Python sender scripts to simulate live traffic. |
| [`tests/`](tests/) | Acceptance and integration checks (`test_api.py`) verifying Auth, Ingestion, and Search. |
| [`docs/`](docs/) | Architecture diagrams, detailed setup guides (Appliance & SaaS), and data flow explanations. |

## Docs
- [Architecture](docs/architecture.md): Components, data flow, normalization, the tenant model, and auth.
- [Appliance Setup](docs/setup_appliance.md): How to run the system locally.
- [SaaS Setup](docs/setup_saas.md): How to deploy to a cloud provider like DigitalOcean.
- [Backend README](backend/README.md): Go backend technical details.
- [Frontend README](frontend/README.md): React SPA technical details.
- [Ingest README](ingest/README.md): Details on the ingestion parsing pipeline.
