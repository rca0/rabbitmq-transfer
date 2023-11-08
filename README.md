# RabbitMQ

A tool to test RabbitMQ connections, and consume/publish/transfer messages between queues. The mode of operation and other parameters are configured using a simple YAML configuration file.

## Consume and Produce

* Connect to RabbitMQ Broker
* Consume messages from a queue
* Produce messages to a queue

## Transfer queue messages

In this mode, all messages that weren't consumed by any client will replicated to another RabbitMQ Server.

## Compare queue size

In this mode, the command will compare the validation of the queue length between RabbitMQ Servers.

## Usage

- Start a local RabbitMQ server

```sh
$ docker-compose up
```

Access the local RabbitMQ Manager: http://localhost:15672

login with `guest` username and password

- Building a binary file

```sh
go build
```

It's necessary write a `config.yaml` file to start the tool

example:
```
---
rabbitmq: 
  amqp_url: amqp://guest:guest@localhost:5672
  queue_name: queue-name-1
```

- Consume some queue

```sh
./rabbitmq --consume
```

- Produce a message in some queue

```sh
./rabbitmq --producer message
```

- Transfer data between RabbitMQ Servers

First step is configure `transfer.yaml` file

* example:
```yaml
servers:
  source: rabbitmq-server-1:5672
  dest: rabbitmq-server-2:5672
vhosts: 
  - name: vhost_name_1
    queues: 
      - name: queue-1
        exchange_letter: queue-4.deadletter
      - name: queue-2
        routing_key: routing-key-value
      - name: queue-3
      - name: queue-4
        routing_key: ""
        exchange_letter: queue-4.deadletter
  - name: vhost_name_2
    queues:
      - name: queue-name-6
      - name: queue-name-5
        routing_key: routing-key-value
      - name: queue-name-7
        exchange_letter: queue-7.deadletter
```

* Running

- Transfer

```sh
./rabbitmq --transfer
```

These settings will connect into the source server and from each *vhost* will copy all messages from the queues to their respectives in the destination server.

- Compare

```sh
./rabbitmq --compare
```

- Diffs

This command will only return the diff of the settings between source and destination servers

 ```sh
./rabbitmq --diff
 ```
