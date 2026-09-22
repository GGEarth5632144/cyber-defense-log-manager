package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "logs.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	DB = db
	createTables()
	seedUsers()
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password_hash TEXT,
		role TEXT,
		tenant TEXT
	);

	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME,
		tenant TEXT,
		source TEXT,
		vendor TEXT,
		product TEXT,
		event_type TEXT,
		event_subtype TEXT,
		severity INTEGER,
		action TEXT,
		src_ip TEXT,
		src_port INTEGER,
		dst_ip TEXT,
		dst_port INTEGER,
		protocol TEXT,
		user TEXT,
		host TEXT,
		process TEXT,
		url TEXT,
		http_method TEXT,
		status_code INTEGER,
		rule_name TEXT,
		rule_id TEXT,
		cloud_account_id TEXT,
		cloud_region TEXT,
		cloud_service TEXT,
		raw JSON,
		tags TEXT,
		received_at DATETIME
	);
	CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp);
	CREATE INDEX IF NOT EXISTS idx_logs_tenant ON logs(tenant);

	CREATE TABLE IF NOT EXISTS alerts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME,
		tenant TEXT,
		message TEXT,
		severity INTEGER,
		is_read BOOLEAN DEFAULT 0
	);
	`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Gracefully add columns to existing logs table in case it was created earlier
	newCols := []string{
		"event_subtype TEXT", "src_port INTEGER", "dst_port INTEGER", "protocol TEXT", "host TEXT", 
		"process TEXT", "url TEXT", "http_method TEXT", "status_code INTEGER", "rule_name TEXT", 
		"rule_id TEXT", "cloud_account_id TEXT", "cloud_region TEXT", "cloud_service TEXT", "tags TEXT",
	}
	for _, col := range newCols {
		DB.Exec(fmt.Sprintf("ALTER TABLE logs ADD COLUMN %s", col))
	}
}

func seedUsers() {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 {
		DB.Exec("INSERT INTO users (username, password_hash, role, tenant) VALUES (?, ?, ?, ?)", "admin", "admin123", "admin", "*")
		DB.Exec("INSERT INTO users (username, password_hash, role, tenant) VALUES (?, ?, ?, ?)", "viewer", "viewer123", "viewer", "demoA")
		DB.Exec("INSERT INTO users (username, password_hash, role, tenant) VALUES (?, ?, ?, ?)", "viewer_b", "viewer123", "viewer", "demoB")
	}
}
