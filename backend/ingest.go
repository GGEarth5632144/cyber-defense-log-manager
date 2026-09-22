package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// HandleHttpIngest accepts JSON logs and normalizes them
func HandleHttpIngest(c *fiber.Ctx) error {
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON payload"})
	}

	rawJSON, _ := json.Marshal(body)

	entry := LogEntry{
		Timestamp:  time.Now(),
		ReceivedAt: time.Now(),
		Raw:        rawJSON,
	}

	// Basic top-level fields mapping
	if val, ok := body["tenant"].(string); ok { entry.Tenant = val } else { entry.Tenant = "default" }
	if val, ok := body["source"].(string); ok { entry.Source = val }
	if val, ok := body["event_type"].(string); ok { entry.EventType = val }
	if val, ok := body["action"].(string); ok { entry.Action = val }
	if val, ok := body["user"].(string); ok { entry.User = val }
	if val, ok := body["ip"].(string); ok { entry.SrcIP = val }
	if val, ok := body["src_ip"].(string); ok { entry.SrcIP = val }
	if val, ok := body["host"].(string); ok { entry.Host = val }
	if val, ok := body["process"].(string); ok { entry.Process = val }
	
	// Float/Int handling for severity and status_code
	if val, ok := body["severity"].(float64); ok { entry.Severity = int(val) }
	if val, ok := body["event_id"].(float64); ok { entry.RuleID = strconv.Itoa(int(val)) }

	// AWS CloudTrail specific mapping
	if cloud, ok := body["cloud"].(map[string]interface{}); ok {
		if val, ok := cloud["service"].(string); ok { entry.CloudService = val }
		if val, ok := cloud["account_id"].(string); ok { entry.CloudAccountID = val }
		if val, ok := cloud["region"].(string); ok { entry.CloudRegion = val }
	}

	// Look for a timestamp
	if tsStr, ok := body["@timestamp"].(string); ok {
		if t, err := time.Parse(time.RFC3339, tsStr); err == nil {
			entry.Timestamp = t
		}
	}

	err := SaveLogEntry(entry)
	if err != nil {
		log.Printf("Error saving log entry: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to ingest log"})
	}

	return c.JSON(fiber.Map{"status": "ok"})
}

func SaveLogEntry(entry LogEntry) error {
	query := `
		INSERT INTO logs (
			timestamp, tenant, source, vendor, product, event_type, event_subtype, severity, action, 
			src_ip, src_port, dst_ip, dst_port, protocol, user, host, process, url, http_method, 
			status_code, rule_name, rule_id, cloud_account_id, cloud_region, cloud_service, raw, tags, received_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := DB.Exec(query,
		entry.Timestamp.UTC().Format("2006-01-02 15:04:05"),
		entry.Tenant,
		entry.Source,
		entry.Vendor,
		entry.Product,
		entry.EventType,
		entry.EventSubtype,
		entry.Severity,
		entry.Action,
		entry.SrcIP,
		entry.SrcPort,
		entry.DstIP,
		entry.DstPort,
		entry.Protocol,
		entry.User,
		entry.Host,
		entry.Process,
		entry.URL,
		entry.HTTPMethod,
		entry.StatusCode,
		entry.RuleName,
		entry.RuleID,
		entry.CloudAccountID,
		entry.CloudRegion,
		entry.CloudService,
		string(entry.Raw),
		entry.Tags,
		entry.ReceivedAt.UTC().Format("2006-01-02 15:04:05"),
	)
	return err
}
