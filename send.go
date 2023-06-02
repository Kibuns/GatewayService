package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func send(username string) {
	queueNames := []string{"delete_user", "delete_twoots", "delete_hashtagtwoots", "delete_credentials"}
    conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/") // locally change rabbitmq to localhost
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    body, err := json.Marshal(username)
    if err != nil {
        log.Panicf("Failed to marshal username: %s", err)
    }

    for _, queueName := range queueNames {
		_, err := ch.QueueDeclare(
			queueName, // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		failOnError(err, "Failed to declare a queue")
        err = ch.PublishWithContext(ctx,
            "",                  // exchange
            queueName,           // routing key
            false,               // mandatory
            false,               // immediate
            amqp.Publishing{
                ContentType: "application/json",
                Body:        body,
            })
        failOnError(err, "Failed to publish a message to queue "+queueName)

        log.Printf(" [x] Sent %s to queue %s\n", body, queueName)
    }
}