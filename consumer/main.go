package consumer

import (
	"fmt"
	"log"
	"strings"

	"github.com/dlpco/devops-tools/rabbitmq/amqp"
	"github.com/dlpco/devops-tools/rabbitmq/loader"
)

func Run(amqpUrl, queueName, vhost string) {
	a, err := amqp.NewConnection(amqpUrl + "/" + vhost)
	if err != nil {
		log.Fatalln(err)
	}
	defer a.Shutdown()

	q := a.QueueDeclare(queueName, loader.Queue{})

	if q.Messages < 1 {
		fmt.Printf("There is no message in the queue: '%s'\n", q.Name)
	}

	cm := a.Consume(q.Name)
	if err != nil {
		log.Fatalln("Consume: ", err)
	}

	for d := range cm.Body {
		log.Printf("['%s'] ['%s'] ['%s'] - Received a message >> '%s'", strings.Split(amqpUrl, "@")[1:], vhost, q.Name, d.Body)
	}
}
