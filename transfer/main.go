package transfer

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/dlpco/devops-tools/rabbitmq/amqp"
	"github.com/dlpco/devops-tools/rabbitmq/loader"
)

func transfer(src, dest, vhost string, values []loader.Queue, wg *sync.WaitGroup) {
	for _, v := range values {
		src, err := amqp.NewConnection(fmt.Sprintf(src + "/" + vhost))
		if err != nil {
			log.Fatalln("Cannot connect to RabbitMQ server", err)
		}
		defer src.Shutdown()

		dst, err := amqp.NewConnection(fmt.Sprintf(dest + "/" + vhost))
		if err != nil {
			log.Fatalln("Cannot connect to RabbitMQ server", err)
		}
		defer dst.Shutdown()

		srcQ := src.QueueDeclare(v.Name, loader.Queue(v))

		if srcQ.Messages < 1 {
			log.Printf("There is no message in queue: '%s'\n", srcQ.Name)
			continue
		}

		c := src.Consume(srcQ.Name)

		dstQ := dst.QueueDeclare(v.Name, loader.Queue(v))

		log.Printf("Consume and publishing in: '%s'\n", dstQ.Name)
	f:
		for {
			select {
			case msg := <-c.Body:
				dst.Publish(dest, vhost, dstQ.Name, msg.Body)
			case <-time.After(2 * time.Second):
				log.Printf("['%s'] - Timeout\n", strings.Split(dest, "@")[1:])
				break f
			}
		}
	}
	wg.Done()
}

func Run() {
	var c loader.Transfer
	cfg := c.GetTransferConfig()

	var wg sync.WaitGroup
	for _, v := range cfg.Vhosts {
		wg.Add(1)
		go transfer(cfg.Servers.Source, cfg.Servers.Dest, v.Name, v.Queues, &wg)
	}
	wg.Wait()
}
