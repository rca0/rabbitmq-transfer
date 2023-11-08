package compare

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/dlpco/devops-tools/rabbitmq/amqp"
	"github.com/dlpco/devops-tools/rabbitmq/loader"
)

func compare(src, dest, vhost string, values []loader.Queue, wg *sync.WaitGroup) {
	for _, v := range values {
		source, err := amqp.NewConnection(fmt.Sprintf(src + "/" + vhost))
		if err != nil {
			log.Fatalln("Cannot connect to RabbitMQ server", err)
		}
		defer source.Shutdown()

		destination, err := amqp.NewConnection(fmt.Sprintf(dest + "/" + vhost))
		if err != nil {
			log.Fatalln("Cannot connect to RabbitMQ server", err)
		}
		defer destination.Shutdown()

		srcQ := source.QueueDeclare(v.Name, loader.Queue(v))
		dstQ := destination.QueueDeclare(v.Name, loader.Queue(v))

		if srcQ.Messages == dstQ.Messages {
			log.Printf("['%s'] - ['%s']: '%#v'\n", strings.Split(src, "@")[1:], strings.Split(dest, "@")[1:], srcQ.Messages)
			continue
		}
		log.Printf("'%s'- ['%s']['%#v'] - ['%s']['%#v']: Queues have different size", srcQ.Name, strings.Split(src, "@")[1:], srcQ.Messages, strings.Split(dest, "@")[1:], dstQ.Messages)
	}
	wg.Done()
}

func Run() {
	var c loader.Transfer
	cfg := c.GetTransferConfig()

	var wg sync.WaitGroup
	for _, v := range cfg.Vhosts {
		wg.Add(1)
		go compare(cfg.Servers.Source, cfg.Servers.Dest, v.Name, v.Queues, &wg)
	}
	wg.Wait()
}
