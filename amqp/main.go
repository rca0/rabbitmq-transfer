package amqp

import (
	"log"
	"strings"

	"github.com/dlpco/devops-tools/rabbitmq/loader"
	"github.com/streadway/amqp"
)

type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type QueueDeclare struct {
	Name     string
	Messages int
}

type Consume struct {
	Name string
	Body <-chan amqp.Delivery
}

func NewConnection(url string) (*Connection, error) {
	var err error
	c := &Connection{}

	log.Printf("Open new connection: %s\n", strings.Split(url, "@")[1:])

	c.conn, err = amqp.Dial(url)
	if err != nil {
		log.Fatalln("Failed to connect to RabbitMQ Server", c)
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		log.Fatalln("Failed to open channel", err)
	}

	return c, nil
}

func (c *Connection) QueueDeclare(queueName string, args loader.Queue) QueueDeclare {
	var attr = make(amqp.Table)

	if args.ExchangeLetter != nil {
		attr["x-dead-letter-exchange"] = *args.ExchangeLetter
	}

	if args.RoutingKey != nil {
		attr["x-dead-letter-routing-key"] = *args.RoutingKey
	}

	q, err := c.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		attr,      // arguments
	)
	if err != nil {
		log.Fatalln("Queue Declare: ", err)
	}
	return QueueDeclare{Name: q.Name, Messages: q.Messages}
}

func (c *Connection) Consume(queueName string) Consume {
	body, err := c.channel.Consume(
		queueName, // queue declared
		"",        // consumer name
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalln("Failed to register a consumer", err)
	}

	return Consume{
		Name: queueName,
		Body: body,
	}
}

func (c *Connection) Publish(host, vhost, queueName string, msgs []byte) {
	err := c.channel.Publish(
		"",        // exchange, empty get the default values
		queueName, // queue declared
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         msgs,
		})

	log.Printf("['%s'] ['%s'] ['%s'] [x] Sent: '%s'", strings.Split(host, "@")[1:], vhost, queueName, msgs)
	if err != nil {
		log.Fatalln("Failed to publish a message", err)
	}
}

func (c *Connection) Shutdown() {
	if err := c.channel.Close(); err != nil {
		log.Fatalln("Close channel failed: %s", err)
	}
	if err := c.conn.Close(); err != nil {
		log.Fatalln("AMQP Connection close error: %s", err)
	}
}
