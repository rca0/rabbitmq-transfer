package main

import (
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dlpco/devops-tools/rabbitmq/compare"
	"github.com/dlpco/devops-tools/rabbitmq/consumer"
	"github.com/dlpco/devops-tools/rabbitmq/diff"
	"github.com/dlpco/devops-tools/rabbitmq/loader"
	"github.com/dlpco/devops-tools/rabbitmq/publish"
	"github.com/dlpco/devops-tools/rabbitmq/transfer"
	"github.com/urfave/cli"
)

func main() {
	logFilename := "rabbitmq-" + time.Now().Format("20060102150405") + ".log"
	f, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	var config loader.Config
	cfg := config.GetConfig()

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "consumer, c",
			Usage: "Consumer message from queue",
		},
		cli.StringFlag{
			Name:  "producer, p",
			Usage: "Produce message in queue",
		},
		cli.BoolFlag{
			Name:  "transfer, t",
			Usage: "Transfer data between rabbitmq servers",
		},
		cli.BoolFlag{
			Name:  "compare, cp",
			Usage: "Compare queues between rabbitmq servers",
		},
		cli.BoolFlag{
			Name:  "diff, df",
			Usage: "Check users/vhosts/exchanges and queues between rabbitmq brokers",
		},
	}
	app.Action = func(c *cli.Context) {
		if c.Bool("consumer") {
			consumer.Run(cfg.Server.AmqpUrl, cfg.Server.QueueName, cfg.Server.Vhost)
		}
		if c.String("producer") != "" {
			if len(os.Args) < 2 {
				log.Fatalln("You must specify message to write in queue")
			}
			msg := strings.Join(os.Args[2:], " ")
			publish.Run(cfg.Server.AmqpUrl, cfg.Server.QueueName, cfg.Server.Vhost, []byte(msg))
		}
		if c.Bool("transfer") {
			transfer.Run()
		}
		if c.Bool("compare") {
			compare.Run()
		}
		if c.Bool("diff") {
			diff.Run()
		}
	}
	app.Run(os.Args)
}
