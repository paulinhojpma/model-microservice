package web

import (
	"fmt"
	"strings"

	"sab.io/escola-service/cache"
	"sab.io/escola-service/database"
	"sab.io/escola-service/logger"
	"sab.io/escola-service/messaging"
	"sab.io/escola-service/storage"
)

// pacote com funções úteis

// GetServices ...
func GetServices(h *Handler) (*database.DataBase, logger.ILogger, cache.ICacheClient, messaging.IMessageClient, storage.IStorage) {
	return h.DB, *h.Logger, *h.Cache, *h.Message, *h.Storage

}

// RowNilInt ...
func rowNilInt(r database.Row, column int) int {
	if r[column] != nil {
		return r.Integer(column)
	}
	return 0
}

func rowNil(r database.Row, column int) string {
	if r[column] != nil {
		return r.String(column)
	}
	return ""
}

func PathMethod(method, resource string) string {
	return strings.ToUpper(fmt.Sprintf("%s:%s", method, resource))
}
