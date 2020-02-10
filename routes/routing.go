package route

import (
	"sab.io/escola-service/messaging"
)

type Manager struct {
	routes map[string]func(*messaging.MessageParam) *messaging.MessageParam
}

// AddRoute ...
func (m *Manager) AddRoute(name string, f func(*messaging.MessageParam) *messaging.MessageParam) {
	m.routes[name] = f
}

// CallService ..
func (m *Manager) CallService(route string, msg *messaging.MessageParam) *messaging.MessageParam {
	newMsg := m.routes[route](msg)
	return newMsg
}

// NewManager ...
func NewManager() *Manager {
	manager := &Manager{}
	manager.routes = make(map[string]func(*messaging.MessageParam) *messaging.MessageParam)
	return manager
}

// ManagerMessage ...
func (m *Manager) ManagerMessage(client messaging.IMessageClient, msg *messaging.MessageParam) error {
	if msg != nil {
		if _, ok := msg.Args["correlationId"]; ok {
			return client.RespondMessage(msg.Args["replyTo"].(string), msg)
		}
		return client.PublishMessage(msg.Args["replyTo"].(string), msg)
	}
	return nil

}
