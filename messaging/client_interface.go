package messaging

import (
	"encoding/json"
)

//IMessageClient ...
type IMessageClient interface {
	connectService(config *OptionsMessageCLient) error
	PublishMessage(routing string, params *MessageParam) error
	ReceiveMessage(routing string) (<-chan MessageParam, error)
	PublishAndReceiveMessage(routing string, params *MessageParam) (*MessageParam, error)
}

type Connection interface {
}

//OptionsMessageCLient opcoes para cofigurar a fila de mensagens
type OptionsMessageCLient struct {
	URL    string                 `json:"url"`
	Driver string                 `json:"driver"`
	Args   map[string]interface{} `json:"args"`
}

//MessageParam ...
type MessageParam struct {
	Method   string                 `json:"method"`
	Query    map[string][]string    `json:"query"`
	Params   map[string]int         `json:"params"`
	Resource string                 `json:"resource"`
	Type     string                 `json:"type"`
	Status   int                    `json:"status"`
	Body     []byte                 `json:"body"`
	Args     map[string]interface{} `json:"args"`
}

//ConfiguraFilaMensagens checa qual a aplicação da fila de mensagens e faz as configuração a partir dele
func (o *OptionsMessageCLient) ConfiguraFilaMensagens() (*IMessageClient, error) {
	var iMessage IMessageClient
	switch o.Driver {
	case "rabbitMQ":

		rab := &Rabbit{}

		errRab := rab.connectService(o)
		if errRab != nil {
			return nil, errRab

		}
		iMessage = rab

	}
	return &iMessage, nil
}

//MountMessageQueue monta a menssagem a ser enviada
func MountMessageQueue(method, resource, typpe string, status int, body *interface{}, params map[string]int, query map[string][]string) (*MessageParam, error) {
	bodyByte, errJSON := json.Marshal(body)
	if errJSON != nil {
		return nil, errJSON
	}
	message := &MessageParam{

		Resource: resource,
		Method:   method,
		Query:    query,
		Params:   params,
		Type:     typpe,
		Status:   status,
		Body:     bodyByte,
	}
	return message, nil
}
