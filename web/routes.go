package web

import (
	"sab.io/escola-service/routes"
)

func NewRoutes(h *Handler) *route.Manager {
	manager := route.NewManager()

	manager.AddRoute("test", h.Test)

	return manager
}
