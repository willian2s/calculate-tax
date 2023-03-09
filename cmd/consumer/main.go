package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/willian2s/calculate-tax/internal/infra/database"
	"github.com/willian2s/calculate-tax/internal/usecase"
	"github.com/willian2s/calculate-tax/pkg/kafka"
	"github.com/willian2s/calculate-tax/pkg/rabbitmq"

	// SQLite3 Driver
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./orders.db")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	repository := database.NewOrderRepository(db)
	usecase := usecase.CalculateFinalPrice{OrderRepository: repository}

	// Kafka
	msgChanKafka := make(chan *ckafka.Message)
	topics := []string{"orders"}
	servers := "host.docker.internal:9094"

	fmt.Println("Kafka consumer has starter")
	go kafka.Consume(topics, servers, msgChanKafka)
	go kafkaWorker(msgChanKafka, usecase)

	// RabbitMQ
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	msgRabbitmqChannel := make(chan amqp.Delivery)

	fmt.Println("RabbitMQ consumer has starter")
	go rabbitmq.Consume(ch, msgRabbitmqChannel)
	rabbitmqWorker(msgRabbitmqChannel, usecase)
}

func kafkaWorker(msgKafka chan *ckafka.Message, uc usecase.CalculateFinalPrice) {
	fmt.Println("Kafka worker has starter")
	for msg := range msgKafka {
		var OrderInputDTO usecase.OrderInputDTO
		err := json.Unmarshal(msg.Value, &OrderInputDTO)
		if err != nil {
			panic(err)
		}

		outputDto, err := uc.Execute(OrderInputDTO)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Kafka has processed order %s\n", outputDto.ID)
	}
}

func rabbitmqWorker(msgChan chan amqp.Delivery, uc usecase.CalculateFinalPrice) {
	fmt.Println("Rabbitmq worker has started")
	for msg := range msgChan {
		var OrderInputDTO usecase.OrderInputDTO
		err := json.Unmarshal(msg.Body, &OrderInputDTO)
		if err != nil {
			panic(err)
		}

		outputDto, err := uc.Execute(OrderInputDTO)
		if err != nil {
			panic(err)
		}

		msg.Ack(false)
		fmt.Printf("Rabbitmq has processed order %s\n", outputDto.ID)
	}
}
