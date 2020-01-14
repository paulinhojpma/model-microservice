package messaging

import (
	"errors"
	"log"
	"time"

	"github.com/streadway/amqp"
)

const (
	//MAXIDDLESCHANNELS numero maximo de channels que podem ser arb
	MAXIDDLESCHANNELS = 20
)

//Rabbit connexao a fila RabbitMQ
type Rabbit struct {
	Connection *amqp.Connection
	Channels   []*Channel
	Exchange   *Exchange `json:"exchange"`
}

//Channel canal de uma connexão e o status
type Channel struct {
	Channel *amqp.Channel
	Used    bool
	Closed  bool
}

//Queue representa um fila de mensagens do serviço RabbitMQ
type Queue struct {
	Status     *amqp.Queue
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

//Exchange ...
type Exchange struct {
	name       string
	kind       string
	durable    bool
	autoDelete bool
	internal   bool
	noWait     bool
	args       amqp.Table
	Binding    map[string][]*Queue
}

//ConnectQueue representa uma conexão com o serviço RabbitMQ
func (rab *Rabbit) connectService(config *OptionsMessageCLient) error {

	//cria uma conexão com o serviço de mensagens
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return err
	}
	rab.Connection = conn

	//gera channels para o pool de channels
	go func() {
		for {
			errChannels := rab.generateIdleChannels(MAXIDDLESCHANNELS)
			if errChannels != nil {
				log.Println("Erro ao gerar novos Channels - ", errChannels)
			}
			time.Sleep(time.Second)
		}

	}()

	channel, errChannel := rab.getActiveChannel()
	if errChannel != nil {
		return errChannel
	}
	//declara os parametros uma exchange exchange com os parametros inseridos
	exchange, errExc := channel.configToExchange(config.Args["exchange"].(map[string]interface{}))
	if errExc != nil {
		return errExc
	}
	rab.Exchange = exchange
	errQueue := channel.configToQueues(config.Args["queue"].([]interface{}), exchange)
	if errQueue != nil {
		return errQueue
	}
	return nil
}

func (rab *Rabbit) generateIdleChannels(maxPoolSize int) error {
	if rab.Connection.IsClosed() {
		return errors.New("Conexão com a fila encerrada")
	}

	if len(rab.Channels) == 0 {
		rab.Channels = make([]*Channel, maxPoolSize)
	}
	//verifica se algumm chaneel está fechado e solicita uma nova conexão
	for i := 0; i < maxPoolSize; i++ {
		if rab.Channels[i] == nil {
			chann, errCha := rab.Connection.Channel()
			if errCha != nil {
				return errCha
			}
			rab.Channels[i].Channel = chann
			rab.Channels[i].Used = false
			rab.Channels[i].Closed = false
		}

	}

	return nil

}

//retorna um canal que esteja ativo e válido
func (rab *Rabbit) getActiveChannel() (*Channel, error) {

	if len(rab.Channels) > 0 {
		for _, channel := range rab.Channels {
			if !channel.Used && !channel.Closed {
				return channel, nil
			}
		}
	}
	return nil, errors.New("Não existe channels disponíveis")
}

//fecha um channel e atualiza o status para encerrado
func (chann *Channel) close() error {
	errChann := chann.Channel.Close()
	if errChann != nil {
		return errChann
	}
	chann.Closed = true
	chann.Used = true
	return nil
}

//PublishMessage publica na fila de mensagens ele manda a mensagem e
//espera a confiramção de ser enviada para a fila, caso seja negado ele envia um erro
func (rab *Rabbit) PublishMessage(routing string, params *MessageParam) error {

	c, erroCh := rab.getActiveChannel()
	if erroCh != nil {
		return erroCh
	}
	c.lockChannel()
	defer c.unLockChannel()
	publishing := &amqp.Publishing{
		ContentType: "application/json",
		Body:        params.Body,
		Headers: amqp.Table{
			"method":   params.Method,
			"params":   params.Params,
			"query":    params.Query,
			"resource": params.Resource,
			"status":   params.Status,
			"type":     params.Type,
		},
	}
	errPubli := rab.publishMessage(routing, c, publishing)
	if errPubli != nil {
		return errPubli
	}
	return nil
}

//ReceiveMessage recebe as mensagens e devolve um channel com as menssagens enviadas assincronicamente
func (rab *Rabbit) ReceiveMessage(routing string) (<-chan MessageParam, error) {
	c, errChannel := rab.getActiveChannel()
	if errChannel != nil {
		return nil, errChannel
	}
	c.lockChannel()
	defer c.unLockChannel()
	msgChan := make(chan MessageParam)
	msg := MessageParam{}
	delivery, err := c.Channel.Consume(
		routing, // queue
		routing, // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return nil, err
	}
	errQ := c.Channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if errQ != nil {
		return nil, errQ
	}
	go func() {
		for d := range delivery {
			msg.Body = d.Body
			msg.Method = d.Headers["method"].(string)
			msg.Params = d.Headers["params"].(map[string]int)
			msg.Query = d.Headers["query"].(map[string][]string)
			msg.Resource = d.Headers["resource"].(string)
			msg.Status = d.Headers["status"].(int)
			msg.Type = d.Headers["type"].(string)
			if d.CorrelationId != "" {
				msg.Args["correlationId"] = d.CorrelationId
				msg.Args["replyTo"] = d.ReplyTo
			}
			msgChan <- msg
			d.Ack(false)
		}
	}()
	time.Sleep(time.Second)
	return msgChan, nil
}

// PublishAndReceiveMessage publica e recebe uma mensagem de forma síncrona.
func (rab *Rabbit) PublishAndReceiveMessage(routing string, params *MessageParam) (*MessageParam, error) {
	c, erroCh := rab.getActiveChannel()
	if erroCh != nil {
		return nil, erroCh
	}
	c.lockChannel()
	defer c.unLockChannel()
	q, errQue := c.Channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if errQue != nil {
		return nil, errQue
	}
	corrID := randomString(32)
	publishing := &amqp.Publishing{
		ContentType:   "application/json",
		Body:          params.Body,
		CorrelationId: corrID,
		ReplyTo:       q.Name,
		Headers: amqp.Table{
			"method":   params.Method,
			"params":   params.Params,
			"query":    params.Query,
			"resource": params.Resource,
			"status":   params.Status,
			"type":     params.Type,
		},
	}
	errPubli := rab.publishMessage(routing, c, publishing)
	if errPubli != nil {
		return nil, errPubli
	}
	msg := MessageParam{}
	delivery, errDelivery := c.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if errDelivery != nil {
		return nil, errDelivery
	}

	for d := range delivery {
		msg.Body = d.Body
		msg.Method = d.Headers["method"].(string)
		msg.Params = d.Headers["params"].(map[string]int)
		msg.Query = d.Headers["query"].(map[string][]string)
		msg.Resource = d.Headers["resource"].(string)
		msg.Status = d.Headers["status"].(int)
		msg.Type = d.Headers["type"].(string)
		if d.CorrelationId != "" {
			msg.Args["correlationId"] = d.CorrelationId
			msg.Args["replyTo"] = d.ReplyTo
		}

		d.Ack(false)
	}
	_, errCloseQue := c.Channel.QueueDelete(q.Name, false, false, false)
	if errCloseQue != nil {
		log.Println("Não foi possível remover a fila")
		c.Closed = true
		go rab.removeCloseChannels()
	}

	return &msg, nil
}

func (rab *Rabbit) publishMessage(routing string, c *Channel, publishing *amqp.Publishing) error {

	confirms := c.Channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	if err := c.Channel.Confirm(false); err != nil {
		log.Fatalf("confirm.select destination: %s", err)
	}
	count := 0

	for {

		errPubli := c.Channel.Publish(
			rab.Exchange.name, // exchange
			routing,           // routing key
			false,             // mandatory
			false,             // immediate
			*publishing)

		if errPubli != nil {
			return errPubli
		}

		count++

		// only ack the source delivery when the destination acks the publishing
		if confirmed := <-confirms; confirmed.Ack {
			log.Println("Mensagem enviada para a fila")
			return nil
		}
		//se após três tentativas a mensagem não for confirmada retorna uma mensagem de erro
		if count > 3 {
			return errors.New("Não foi possível entregar a mensagem")
		}

	}

}
func (c *Channel) configToExchange(exc map[string]interface{}) (*Exchange, error) {
	exchange := &Exchange{}

	if exc["autoDelete"] == nil {
		return nil, errors.New("Valor de configuração 'autoDelete' vazio")
	}
	if exc["durable"] == nil {
		return nil, errors.New("Valor de configuração 'durable' vazio")
	}
	if exc["internal"] == nil {
		return nil, errors.New("Valor de configuração 'internal' vazio")
	}
	if exc["kind"] == nil {
		return nil, errors.New("Valor de configuração 'kind' vazio")
	}
	if exc["name"] == nil {
		return nil, errors.New("Valor de configuração 'name' vazio")
	}
	if exc["noWait"] == nil {
		return nil, errors.New("Valor de configuração 'noWait' vazio")
	}

	exchange.autoDelete = exc["autoDelete"].(bool)
	exchange.durable = exc["durable"].(bool)
	exchange.internal = exc["internal"].(bool)
	exchange.kind = exc["kind"].(string)
	exchange.name = exc["name"].(string)
	exchange.noWait = exc["noWait"].(bool)
	errChannel := c.Channel.ExchangeDeclare(
		exchange.name,       // name
		"direct",            // type
		exchange.durable,    // durable
		exchange.autoDelete, // auto-deleted
		exchange.internal,   // internal
		exchange.noWait,     // no-wait
		nil,                 // arguments
	)
	if errChannel != nil {
		return nil, errChannel
	}
	// err := json.Unmarshal(exc, exchange)
	// if err != nil {
	// 	return nil, err
	// }
	return exchange, nil
}

//configToQueues declara as queues e associa a uma exchange com uma bidingQueue
func (c *Channel) configToQueues(qArgs []interface{}, exchange *Exchange) error {
	queues := make([]*Queue, len(qArgs))
	bindingQueues := make(map[string][]*Queue)
	for index, queueInt := range qArgs {
		queue := queueInt.(map[string]interface{})
		name := queue["name"].(string)
		queues[index] = &Queue{}
		queues[index].Name = name
		queues[index].AutoDelete = queue["autoDelete"].(bool)
		queues[index].Durable = queue["durable"].(bool)
		queues[index].Exclusive = queue["exclusive"].(bool)
		queues[index].NoWait = queue["noWait"].(bool)

		queueRab, errChannel := c.Channel.QueueDeclare(
			queues[index].Name,       // name
			queues[index].Durable,    // durable
			queues[index].AutoDelete, // delete when unused
			queues[index].Exclusive,  // exclusive
			queues[index].NoWait,     // no-wait
			nil,                      // arguments
		)
		queues[index].Status = &queueRab
		if errChannel != nil {
			return errChannel
		}
		if len(queue["bindingKeys"].([]interface{})) > 0 {
			for _, binding := range queue["bindingKeys"].([]interface{}) {

				c.Channel.QueueBind(
					queues[index].Name, // queue name
					binding.(string),   // routing key
					exchange.name,      // exchange
					false,
					nil)
				bindingQueues[binding.(string)] = append(bindingQueues[binding.(string)], queues[index])
			}
		}
	}

	exchange.Binding = bindingQueues
	return nil
}

//lockChannel trava o canal atual para não ser utilizado por outro serviço
func (c *Channel) lockChannel() {
	c.Used = true
}

//unLockChannel destrava o canal atual que estava sendo utilizado
func (c *Channel) unLockChannel() {
	c.Used = false
}

func (rab *Rabbit) removeCloseChannels() {
	for i, channel := range rab.Channels {
		if channel.Closed {
			rab.Channels[i] = nil
		}
	}
}
