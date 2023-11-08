package loader

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func (c *Config) GetConfig() *Config {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalln("yaml file error %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalln("unmarshal error %v", err)
	}
	return c
}

func (t *Transfer) GetTransferConfig() *Transfer {
	yamlFile, err := ioutil.ReadFile("transfer.yaml")
	if err != nil {
		log.Fatalln("yaml file error", err)
	}

	if err = yaml.Unmarshal(yamlFile, &t); err != nil {
		log.Fatalln("unmarshal error", err)
	}
	return t
}
