package main

import (
	"log"
	"time"
)

func startRetentionEngine() {
	log.Println("Retention engine started (7-day policy)")
	for {
		// Run cleanup every hour
		time.Sleep(1 * time.Hour)
		
		// Delete logs older than 7 days
		res, err := DB.Exec("DELETE FROM logs WHERE timestamp < datetime('now', '-7 days')")
		if err != nil {
			log.Printf("Retention error: %v", err)
			continue
		}
		
		rows, _ := res.RowsAffected()
		if rows > 0 {
			log.Printf("Retention cleanup: deleted %d old log entries", rows)
		}
	}
}
