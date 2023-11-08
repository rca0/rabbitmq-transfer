package publish

import (
	"log"

	"github.com/dlpco/devops-tools/rabbitmq/amqp"
	"github.com/dlpco/devops-tools/rabbitmq/loader"
)

func Run(amqpUrl, queueName, vhost string, msgs []byte) {
	a, err := amqp.NewConnection(amqpUrl + "/" + vhost)
	if err != nil {
		log.Fatalln(err)
	}
	defer a.Shutdown()

	q := a.QueueDeclare(queueName, loader.Queue{})

	a.Publish(amqpUrl, vhost, q.Name, msgs)
}
