package route

import (
	"sab.io/escola-service/messaging"
)

type Manager struct {
	Routes map[string]func(*messaging.MessageParam) *messaging.MessageParam
}

// AddRoute ...
func (m *Manager) AddRoute(name string, f func(*messaging.MessageParam) *messaging.MessageParam) {
	m.Routes[name] = f
}

// CallService ..
func (m *Manager) CallService(route string, msg *messaging.MessageParam) *messaging.MessageParam {
	newMsg := m.Routes[route](msg)
	return newMsg
}

// NewManager ...
func NewManager() *Manager {
	manager := &Manager{}
	manager.Routes = make(map[string]func(*messaging.MessageParam) *messaging.MessageParam)
	return manager
}
