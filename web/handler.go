package web

import (
	"log"

	"sab.io/escola-service/cache"
	"sab.io/escola-service/database"
	"sab.io/escola-service/logger"
	"sab.io/escola-service/messaging"
	"sab.io/escola-service/storage"
)

// Handler ...
type Handler struct {
	// Relic        newrelic.Application
	// ClientRedis  *redis.Client
	DB *database.DataBase
	// HostsVálidos []string
	// Log          *snetlog.Log
	// GodinAuthURL string
	// Upgrader     websocket.Upgrader
	// EmailConf    map[string]string
	Message *messaging.IMessageClient
	Cache   *cache.ICacheClient
	Logger  *logger.ILogger
	Storage *storage.IStorage
}

func (h *Handler) Test(m *messaging.MessageParam) *messaging.MessageParam {
	log.Println("Mensagem recebida")
	return m
}

func (h *Handler) CadastrarEscola(m *messaging.MessageParam) *messaging.MessageParam {

	return nil
}
