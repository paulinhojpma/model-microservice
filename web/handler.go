package web

import "github.com/sab.io/escola-service/messaging"

type Handler struct {
	// Relic        newrelic.Application
	// ClientRedis  *redis.Client
	// DB           *database.DataBase
	// HostsVálidos []string
	// Log          *snetlog.Log
	// GodinAuthURL string
	// Upgrader     websocket.Upgrader
	// EmailConf    map[string]string
	Messaging *messaging.IMessageClient
}
