# Data Engineering Complete Pipeline for Ride-Sharing Company Telemetry

## Overview

This project implements a complete, production-grade data engineering pipeline for ingesting, processing, and storing real-time telemetry data from a ride-sharing company's vehicles. The pipeline leverages Apache Kafka for streaming, PostgreSQL for storage, and is built using Go for both the producer and consumer microservices. The system is fully containerized using Docker and orchestrated with Docker Compose for easy local development and testing.

## Architecture

- **Producer Service (Go + Fiber):**
  - Exposes an HTTP API endpoint to receive telemetry data (latitude, longitude, speed, timestamp, trip ID).
  - Publishes incoming telemetry data to a Kafka topic.
- **Kafka + Zookeeper:**
  - Kafka acts as the message broker for real-time streaming of telemetry data.
  - Zookeeper manages Kafka's distributed state.
- **Consumer Service (Go):**
  - Subscribes to the Kafka topic, consumes telemetry messages, and writes them to a PostgreSQL database.
- **PostgreSQL:**
  - Stores all ingested telemetry data for further analysis and reporting.
- **Data Simulation Script (Python):**
  - Simulates vehicle telemetry and sends it to the producer's HTTP API for end-to-end testing.

```
[Vehicle Telemetry] → [Producer API] → [Kafka] → [Consumer] → [PostgreSQL]
```

## Prerequisites

- [Docker](https://www.docker.com/get-started) & [Docker Compose](https://docs.docker.com/compose/)
- [Go 1.22+](https://golang.org/dl/) (for local development, not needed for Dockerized run)
- [Python 3.7+](https://www.python.org/downloads/) (for running the simulation script)

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/kotappa19/Data-Engineering-Complete-Pipeline-for-Ride-sharing-company-telemetry.git
cd Data-Engineering-Complete-Pipeline-for-Ride-sharing-company-telemetry/go
```

### 2. Create a `.env` File

Create a `.env` file in the `go/` directory with the following content (adjust as needed):

```
# Kafka
BOOTSTRAP_SERVER=kafka:9092
KAFKA_TOPIC=telemetry

# PostgreSQL
DB_HOST=postgres
DB_PORT=5432
DB_USER=replaceme
DB_PASS=replaceme
DB_NAME=telemetry_db
DB_SSLMODE=disable
```

### 3. Build and Start the Pipeline

Use Docker Compose to build and start all services:

```bash
docker-compose up --build
```

This will start:
- Zookeeper (port 2181)
- Kafka (port 9092)
- PostgreSQL (port 5432)
- Producer API (port 8080)
- Consumer

### 4. Verify the Producer API

Check the health endpoint:

```bash
curl http://localhost:8080/health
```
You should see `OK` if the producer is running.

### 5. Simulate Telemetry Data (Optional)

In a new terminal, from the project root:

```bash
cd scripts
python3 -m venv venv           # Create a virtual environment (recommended)
source venv/bin/activate       # Activate the virtual environment
pip3 install requests           # Install dependencies inside the venv
python3 simulate_data.py        # Run the simulation script
```

This will send simulated telemetry data to the producer API, which will flow through Kafka and be stored in PostgreSQL.

### 6. Accessing the Data

You can connect to the PostgreSQL database using any client:
- Host: `localhost`
- Port: `5432`
- User: `postgres`
- Password: `password`
- Database: `telemetry_db`

Example (using psql):
```bash
psql -h localhost -U postgres -d telemetry_db
```

## Project Structure

```
├── go/
│   ├── producer/         # Producer service (Go + Fiber)
│   ├── consumer/         # Consumer service (Go)
│   ├── models/           # Data models
│   ├── storage/          # Database connection logic
│   ├── Dockerfile-*      # Dockerfiles for producer/consumer
│   ├── docker-compose.yml
│   └── wait-for-it.sh    # Wait script for service dependencies
├── scripts/
│   └── simulate_data.py  # Python script to simulate telemetry
└── Kafka/                # Kafka config/data (optional)
```

## Useful Commands

- **Stop all services:**
  ```bash
  docker-compose down
  ```
- **Rebuild services:**
  ```bash
  docker-compose up --build
  ```
- **View logs:**
  ```bash
  docker-compose logs -f
  ```

## Notes
- All services are containerized; no local Go or PostgreSQL installation is required for running the pipeline.
- The simulation script is optional but recommended for testing the end-to-end flow.
- For production, consider securing environment variables and using managed Kafka/Postgres services.
