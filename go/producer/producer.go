package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Telemetry struct {
	TripId    string  `json:"trip_id"`
	Lat       float64 `json:"lat"`
	Long      float64 `json:"long"`
	Speed     float64 `json:"speed"`
	Timestamp string  `json:"timestamp"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error in loading env file: ", err)
		return
	}

	bootstrapServer := os.Getenv("BOOTSTRAP_SERVER")
	if bootstrapServer == "" {
		fmt.Println("BOOTSTRAP_SERVER env variable not set")
		return
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		fmt.Println("KAFKA_TOPIC env variable not set")
		return
	}

	// Create a new Kafka producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	// Create a new Fiber app
	app := fiber.New()

	app.Post("/telemetry", func(c *fiber.Ctx) error {
		// Parse the request body into a Telemetry struct
		telemetry := Telemetry{}
		err := c.BodyParser(&telemetry)
		if err != nil {
			c.Status(http.StatusUnprocessableEntity).JSON(
				&fiber.Map{"message": "request failed"})
			return err
		}

		// Serialize the Telemetry struct to JSON
		telemetryJSON, err := json.Marshal(telemetry)
		if err != nil {
			log.Printf("Failed to marshal telemetry: %s", err)
			c.Status(http.StatusInternalServerError).JSON(
				&fiber.Map{"message": "Failed to process telemetry"})
			return err
		}

		// Produce the message to the Kafka topic
		producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &kafkaTopic,
				Partition: kafka.PartitionAny,
			},
			Value: telemetryJSON,
		}, nil)

		// Return a success message
		c.Status(http.StatusOK).JSON(
			&fiber.Map{"message": "Telemetry sent to kafka topic " + kafkaTopic + " successfully"})

		return nil
	})

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Listen(":8080")
}
