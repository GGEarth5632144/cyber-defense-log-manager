package main

import (
	"encoding/json"
	"time"
)

type User struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
	Tenant       string `json:"tenant"`
}

type LogEntry struct {
	ID             uint            `json:"id"`
	Timestamp      time.Time       `json:"@timestamp"`
	Tenant         string          `json:"tenant"`
	Source         string          `json:"source"` // firewall|crowdstrike|aws|m365|ad|api|network
	Vendor         string          `json:"vendor"`
	Product        string          `json:"product"`
	EventType      string          `json:"event_type"`
	EventSubtype   string          `json:"event_subtype"`
	Severity       int             `json:"severity"`
	Action         string          `json:"action"` // allow|deny|create|delete|login|logout|alert
	SrcIP          string          `json:"src_ip"`
	SrcPort        int             `json:"src_port"`
	DstIP          string          `json:"dst_ip"`
	DstPort        int             `json:"dst_port"`
	Protocol       string          `json:"protocol"`
	User           string          `json:"user"`
	Host           string          `json:"host"`
	Process        string          `json:"process"`
	URL            string          `json:"url"`
	HTTPMethod     string          `json:"http_method"`
	StatusCode     int             `json:"status_code"`
	RuleName       string          `json:"rule_name"`
	RuleID         string          `json:"rule_id"`
	CloudAccountID string          `json:"cloud_account_id"`
	CloudRegion    string          `json:"cloud_region"`
	CloudService   string          `json:"cloud_service"`
	Raw            json.RawMessage `json:"raw"`
	Tags           string          `json:"_tags"` // Stored as JSON string array
	ReceivedAt     time.Time       `json:"received_at"`
}

type DashboardStats struct {
	TotalLogs     int64            `json:"total_logs"`
	TopIPs        []TopStringCount `json:"top_ips"`
	TopUsers      []TopStringCount `json:"top_users"`
	TopEventTypes []TopStringCount `json:"top_event_types"`
}

type TopStringCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type TimeSeriesPoint struct {
	Time  string `json:"time"`
	Count int    `json:"count"`
}

type Alert struct {
	ID        uint      `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Tenant    string    `json:"tenant"`
	Message   string    `json:"message"`
	Severity  int       `json:"severity"`
	IsRead    bool      `json:"is_read"`
}
