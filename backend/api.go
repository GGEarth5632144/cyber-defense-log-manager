package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Helper to append tenant and time filters to queries
func getTenantFilter(c *fiber.Ctx) (string, []interface{}) {
	tenant := c.Locals("tenant").(string)
	requestedTenant := c.Query("tenant", "")
	timeRange := c.Query("time", "24h")
	
	var timeMod string
	if timeRange == "7d" {
		timeMod = "-7 days"
	} else if timeRange == "1h" {
		timeMod = "-1 hours"
	} else {
		timeMod = "-24 hours"
	}
	
	filter := "timestamp > datetime('now', ?)"
	args := []interface{}{timeMod}
	
	if tenant == "*" {
		if requestedTenant != "" && requestedTenant != "all" {
			filter += " AND tenant = ?"
			args = append(args, requestedTenant)
		}
	} else {
		filter += " AND tenant = ?"
		args = append(args, tenant)
	}
	
	return filter, args
}

func GetDashboardStats(c *fiber.Ctx) error {
	filter, args := getTenantFilter(c)

	var stats DashboardStats

	DB.QueryRow("SELECT COUNT(*) FROM logs WHERE "+filter, args...).Scan(&stats.TotalLogs)

	rows, err := DB.Query("SELECT src_ip, COUNT(*) as c FROM logs WHERE src_ip IS NOT NULL AND src_ip != '' AND "+filter+" GROUP BY src_ip ORDER BY c DESC LIMIT 5", args...)
	if err == nil {
		for rows.Next() {
			var ip string
			var count int
			rows.Scan(&ip, &count)
			stats.TopIPs = append(stats.TopIPs, TopStringCount{Name: ip, Count: count})
		}
		rows.Close()
	}

	rows, err = DB.Query("SELECT user, COUNT(*) as c FROM logs WHERE user IS NOT NULL AND user != '' AND "+filter+" GROUP BY user ORDER BY c DESC LIMIT 5", args...)
	if err == nil {
		for rows.Next() {
			var u string
			var count int
			rows.Scan(&u, &count)
			stats.TopUsers = append(stats.TopUsers, TopStringCount{Name: u, Count: count})
		}
		rows.Close()
	}

	rows, err = DB.Query("SELECT event_type, COUNT(*) as c FROM logs WHERE event_type IS NOT NULL AND event_type != '' AND "+filter+" GROUP BY event_type ORDER BY c DESC LIMIT 5", args...)
	if err == nil {
		for rows.Next() {
			var evt string
			var count int
			rows.Scan(&evt, &count)
			stats.TopEventTypes = append(stats.TopEventTypes, TopStringCount{Name: evt, Count: count})
		}
		rows.Close()
	}

	return c.JSON(stats)
}

func GetTimeline(c *fiber.Ctx) error {
	filter, args := getTenantFilter(c)
	
	query := "SELECT strftime('%Y-%m-%d %H:00:00', timestamp) as hr, COUNT(*) as c FROM logs WHERE " + filter + " GROUP BY hr ORDER BY hr ASC LIMIT 24"
	
	rows, err := DB.Query(query, args...)
	var timeline []TimeSeriesPoint
	if err != nil {
		log.Printf("Timeline query error: %v", err)
		return c.JSON(timeline)
	}
	defer rows.Close()
	if rows.Err() != nil {
		log.Printf("Timeline rows error: %v", rows.Err())
		return c.JSON(timeline)
	}
	for rows.Next() {
		var hr string
		var count int
		rows.Scan(&hr, &count)
		timeline = append(timeline, TimeSeriesPoint{Time: hr, Count: count})
	}

	return c.JSON(timeline)
}

func GetLogs(c *fiber.Ctx) error {
	filter, args := getTenantFilter(c)
	
	rows, err := DB.Query("SELECT id, timestamp, tenant, source, event_type, action, src_ip, user, raw FROM logs WHERE "+filter+" ORDER BY timestamp DESC LIMIT 50", args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch logs"})
	}
	defer rows.Close()

	var logs []map[string]interface{}
	if rows.Err() != nil {
		log.Printf("Logs rows error: %v", rows.Err())
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch logs"})
	}
	for rows.Next() {
		var id int
		var ts time.Time
		var tenant, source, eventType, action, srcIp, user, raw string
		
		err := rows.Scan(&id, &ts, &tenant, &source, &eventType, &action, &srcIp, &user, &raw)
		if err == nil {
			logs = append(logs, map[string]interface{}{
				"id": id,
				"timestamp": ts.Format(time.RFC3339),
				"tenant": tenant,
				"source": source,
				"event_type": eventType,
				"action": action,
				"src_ip": srcIp,
				"user": user,
				"raw": raw,
			})
		}
	}
	
	if logs == nil {
		logs = make([]map[string]interface{}, 0)
	}

	return c.JSON(logs)
}

func GetAlerts(c *fiber.Ctx) error {
	filter, args := getTenantFilter(c)
	
	rows, err := DB.Query("SELECT id, timestamp, tenant, message, severity FROM alerts WHERE "+filter+" ORDER BY timestamp DESC LIMIT 10", args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch alerts"})
	}
	defer rows.Close()

	var alerts []Alert
	if rows.Err() != nil {
		log.Printf("Alerts rows error: %v", rows.Err())
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch alerts"})
	}
	for rows.Next() {
		var a Alert
		rows.Scan(&a.ID, &a.Timestamp, &a.Tenant, &a.Message, &a.Severity)
		alerts = append(alerts, a)
	}
	
	if alerts == nil {
		alerts = make([]Alert, 0)
	}

	return c.JSON(alerts)
}
