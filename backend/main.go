package main

import (
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	InitDB()

	go startSyslogServer()
	go startAlertingEngine()
	go startRetentionEngine() // 7-day retention cleanup

	app := fiber.New(fiber.Config{ DisableStartupMessage: true })
	app.Use(cors.New(cors.Config{ AllowOrigins: "*" }))
	app.Use(logger.New())

	app.Post("/login", LoginHandler)
	app.Post("/ingest", HandleHttpIngest)

	api := app.Group("/api", AuthMiddleware)
	api.Get("/stats", GetDashboardStats)
	api.Get("/timeline", GetTimeline)
	api.Get("/logs", GetLogs)
	api.Get("/alerts", GetAlerts)

	port := "8080"
	log.Printf("Starting HTTP server on :%s", port)
	app.Listen(":" + port)
}

func startSyslogServer() {
	addr, _ := net.ResolveUDPAddr("udp", ":514")
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Printf("Failed to start Syslog server: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Listening for Syslog on UDP :514")
	buffer := make([]byte, 2048)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil { continue }
		handleSyslogMessage(string(buffer[:n]))
	}
}

func handleSyslogMessage(msg string) {
	entry := LogEntry{
		Timestamp:  time.Now(),
		Tenant:     "demoA", // Default for direct syslog
		Source:     "syslog",
		Raw:        []byte(`{"raw": "` + strings.ReplaceAll(strings.TrimSpace(msg), `"`, `\"`) + `"}`),
		ReceivedAt: time.Now(),
	}

	// Extremely simplified kv parsing for the demo string formats
	parts := strings.Split(msg, " ")
	for _, p := range parts {
		if strings.Contains(p, "=") {
			kv := strings.SplitN(p, "=", 2)
			if len(kv) != 2 { continue }
			k, v := kv[0], kv[1]
			
			switch k {
			case "action": entry.Action = v
			case "src": entry.SrcIP = v
			case "dst": entry.DstIP = v
			case "user": entry.User = v
			case "vendor": entry.Vendor = v
			case "product": entry.Product = v
			case "proto": entry.Protocol = v
			case "spt": 
				if p, err := strconv.Atoi(v); err == nil { entry.SrcPort = p }
			case "dpt":
				if p, err := strconv.Atoi(v); err == nil { entry.DstPort = p }
			case "event": entry.EventType = v
			case "mac": entry.Host = v // using host for mac in network logs
			}
		}
	}

	// Router vs Firewall heuristic
	if strings.Contains(msg, "event=link-down") {
		entry.Source = "network"
		entry.Vendor = "Router"
	} else if entry.Vendor == "demo" && entry.Product == "ngfw" {
		entry.Source = "firewall"
	}

	SaveLogEntry(entry)
}
