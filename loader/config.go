package loader

type Config struct {
	Server struct {
		AmqpUrl   string `yaml:"amqp_url"`
		Vhost     string `yaml:"vhost_name"`
		QueueName string `yaml:"queue_name"`
	} `yaml:"rabbitmq"`
}
