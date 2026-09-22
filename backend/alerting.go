package main

import (
	"log"
	"time"
)

type AlertCondition struct {
	Tenant string
	SrcIP  string
	Count  int
}

func startAlertingEngine() {
	// Simple polling engine that checks for repeated failed logins
	log.Println("Alerting engine started")
	
	for {
		time.Sleep(1 * time.Minute)
		
		// Look for >=3 denied actions from the same IP in the last 5 minutes
		// Using pure SQLite datetime functions avoids Go timezone serialization issues
		query := `
			SELECT tenant, src_ip, COUNT(*) as c 
			FROM logs 
			WHERE action = 'deny' AND timestamp > datetime('now', '-5 minutes') AND src_ip IS NOT NULL AND src_ip != ''
			GROUP BY tenant, src_ip 
			HAVING c >= 3
		`
		
		rows, err := DB.Query(query)
		if err != nil {
			log.Printf("Alerting query error: %v", err)
			continue
		}
		
		if rows.Err() != nil {
			log.Printf("Alerting rows error: %v", rows.Err())
			rows.Close()
			continue
		}
		var conditions []AlertCondition
		for rows.Next() {
			var cond AlertCondition
			rows.Scan(&cond.Tenant, &cond.SrcIP, &cond.Count)
			conditions = append(conditions, cond)
		}
		rows.Close() // CRITICAL: Close rows before writing to avoid SQLite database locks
		
		for _, cond := range conditions {
			var recentAlertCount int
			err := DB.QueryRow("SELECT COUNT(*) FROM alerts WHERE tenant = ? AND message LIKE ? AND timestamp > datetime('now', '-5 minutes')", cond.Tenant, "%"+cond.SrcIP+"%").Scan(&recentAlertCount)
			
			if err == nil && recentAlertCount == 0 {
				msg := "Repeated failed actions/logins detected from IP: " + cond.SrcIP
				_, err = DB.Exec("INSERT INTO alerts (timestamp, tenant, message, severity) VALUES (datetime('now'), ?, ?, ?)", cond.Tenant, msg, 8)
				if err != nil {
					log.Printf("Failed to insert alert: %v", err)
				} else {
					log.Printf("ALERT TRIGGERED: %s", msg)
				}
			}
		}
	}
}
