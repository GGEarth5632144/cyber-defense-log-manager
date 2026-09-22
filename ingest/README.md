# Ingestion Architecture

This directory serves as the documentation hub for the Log Management System's ingestion pipeline. 

> **Note on Code Location:** To guarantee database stability and prevent SQLite file-locking issues (`database is locked` errors), the actual ingestion code is tightly coupled within the Go backend service rather than running as a standalone microservice in this folder. The ingestion logic can be found in `/backend/ingest.go` and `/backend/main.go`.

## Supported Ingestion Protocols
The system natively supports two real-time ingestion protocols, fulfilling the multi-source requirement:

### 1. HTTP JSON API (Port 8080)
- **Endpoint:** `POST /ingest`
- **Location:** `/backend/ingest.go`
- **Behavior:** Accepts raw JSON payloads. It maps deeply nested fields (such as AWS CloudTrail's `cloud.account_id`) and standardizes them into the common `LogEntry` schema. 
- **Usage:** Ideal for CrowdStrike, Microsoft 365, AWS CloudTrail, and Microsoft AD logs.

### 2. UDP Syslog (Port 514)
- **Endpoint:** `UDP :514`
- **Location:** `/backend/main.go` (`startSyslogServer`)
- **Behavior:** Runs a dedicated background goroutine listening for raw UDP packets. It uses heuristic parsing to extract key-value pairs (e.g., `src=10.0.1.10`, `dpt=53`) and automatically tags the `source` as either `firewall` or `network` based on the payload signature.

## The Normalization Pipeline
Regardless of which protocol receives the data, all logs undergo normalization before hitting the storage layer. 

1. **Timestamp Standardization:** Time formats are parsed and strictly converted to UTC `YYYY-MM-DD HH:MM:SS` strings to guarantee accurate range queries in SQLite.
2. **Schema Mapping:** Source-specific fields (like `mac` in a network log or `ip` in a login log) are mapped to standard columns (`host` and `src_ip`).
3. **Raw Preservation:** The original, unadulterated message payload is always preserved in the `raw` JSON column for auditing purposes.
