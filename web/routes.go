package web

import (
	route "sab.io/escola-service/routes"
)

func NewRoutes(h *Handler) *route.Manager {
	manager := route.NewManager()

	manager.AddRoute("test", h.Test)
	manager.AddRoute("POST:ESCOLA", h.CadastrarEscola)

	return manager
}
