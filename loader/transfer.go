package loader

type Transfer struct {
	Servers struct {
		Source string `yaml:"source"`
		Dest   string `yaml:"dest"`
	} `yaml:"servers"`
	Vhosts []struct {
		Name       string   `yaml:"name"`
		Deadletter []string `yaml:"deadletter"`
		Queues     []Queue  `yaml:"queues"`
	} `yaml:"vhosts"`
}

type Queue struct {
	Name           string  `yaml:"name"`
	ExchangeLetter *string `yaml:"exchange_letter,omitempty"`
	RoutingKey     *string `yaml:"routing_key,omitempty"`
}
