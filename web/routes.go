package web

import (
	route "sab.io/escola-service/routes"
)

// NewRoutes ...
func NewRoutes(h *Handler) *route.Manager {
	manager := route.NewManager()

	manager.AddRoute("POST:ESCOLA", h.CadastrarEscola)
	manager.AddRoute("GET:ESCOLAS", h.GetEscolas)
	manager.AddRoute("GET:ESCOLA", h.GetEscola)
	manager.AddRoute("POST:UNIDADE", h.CadastrarUnidade)
	manager.AddRoute("POST:DISCIPLINA", h.CadastrarDisciplina)
	manager.AddRoute("GET:DISCIPLINAS", h.GetDisciplinas)
	manager.AddRoute("GET:DISCIPLINA", h.GetDisciplina)
	manager.AddRoute("PUT:DISCIPLINA", h.AtualizaDisciplina)
	manager.AddRoute("DEL:DISCIPLINA", h.DeleteDisciplina)
	manager.AddRoute("DEL:EMENTA", h.DeleteEmenta)

	return manager
}
