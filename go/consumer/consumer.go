package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
	"github.com/kotappa19/Data-Engineering-Complete-Pipeline-for-Ride-sharing-company-telemetry.git/models"
	"github.com/kotappa19/Data-Engineering-Complete-Pipeline-for-Ride-sharing-company-telemetry.git/storage"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error in loading env file: ", err)
		return
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not load database")
		return
	}

	err = models.MigrateTelemetry(db)
	if err != nil {
		log.Fatal("Could not migrate db")
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

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		// User-specific properties that you must set
		"bootstrap.servers": bootstrapServer,

		// Fixed properties
		"group.id":          "vehicle-group",
		"auto.offset.reset": "earliest"})

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	_ = consumer.SubscribeTopics([]string{kafkaTopic}, nil)
	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			telemetry := models.Telemetry{}
			err = json.Unmarshal(ev.Value, &telemetry)
			if err != nil {
				log.Printf("Failed to unmarshal telemetry: %s", err)
				continue
			}
			err = db.Create(&telemetry).Error
			if err != nil {
				log.Printf("Failed to save telemetry to database: %s", err)
				continue
			}
			fmt.Println("Telemetry saved to database successfully")
		}
	}

	consumer.Close()
}
