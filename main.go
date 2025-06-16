package main

import (
	"context"
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	"log"
	"log_processor/internal/llm"
	"log_processor/internal/pubsub"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	amqpConnectionURL string
	dbDSN             string
	llm               struct {
		endpoint      string
		authorization string
	}
	messagingQueue struct {
		exchange   string
		kind       string
		queue      string
		routingKey string
		durable    bool
	}
}

func main() {
	var cfg config
	setUpConfig(&cfg)
	validateConfig(&cfg)

	ps, err := pubsub.NewPubSubConnection(cfg.amqpConnectionURL)
	if err != nil {
		log.Fatalf("Failed to create new pubsub connection. Error: %v", err)
	}

	llmConfig := &llm.LLM{
		URL:           cfg.llm.endpoint,
		Authorization: cfg.llm.authorization,
	}

	db, err := connectDB(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = ps.ConsumeMessage(cfg.messagingQueue.exchange, cfg.messagingQueue.kind, cfg.messagingQueue.queue, cfg.messagingQueue.routingKey, cfg.messagingQueue.durable, llmConfig, db)
	if err != nil {
		log.Fatalf("Failed to consume messages. Error: %v", err)
	}

	log.Println("Running llm log processor service...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	stopSignal := <-stop

	if stopSignal != nil {
		log.Printf("Detected %v signal. Shutting down llm log processor service \n", <-stop)
		os.Exit(0)
	}
}

func connectDB(cfg *config) (*sql.DB, error) {
	log.Println("Attempting to connect to DB with:", cfg.dbDSN)
	db, err := sql.Open("postgres", cfg.dbDSN)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setUpConfig(cfg *config) {
	flag.StringVar(&cfg.amqpConnectionURL, "amqp-connection-url", "", "amqp connection url")
	flag.StringVar(&cfg.llm.endpoint, "llm-endpoint", "", "llm url")
	flag.StringVar(&cfg.llm.authorization, "llm-auth", "", "llm auth key")
	flag.StringVar(&cfg.dbDSN, "db-dsn", "", "db dsn")
	flag.StringVar(&cfg.messagingQueue.exchange, "mq-exchange", "", "rabbitmq exchange name")
	flag.StringVar(&cfg.messagingQueue.kind, "mq-exchange-kind", "", "rabbitmq exchange kind")
	flag.StringVar(&cfg.messagingQueue.queue, "mq-queue", "", "rabbitmq queue name")
	flag.StringVar(&cfg.messagingQueue.routingKey, "mq-routing-key", "", "rabbitmq routing key")
	flag.BoolVar(&cfg.messagingQueue.durable, "mq-durable", true, "rabbitmq queue durability ")

	flag.Parse()

	// Fallbacks from environment
	if cfg.llm.endpoint == "" {
		cfg.llm.endpoint = os.Getenv("LLM_ENDPOINT")
	}
	if cfg.llm.authorization == "" {
		cfg.llm.authorization = os.Getenv("LLM_AUTHORIZATION_KEY")
	}
	if cfg.dbDSN == "" {
		cfg.dbDSN = os.Getenv("DB_DSN")
	}
	if cfg.messagingQueue.exchange == "" {
		cfg.messagingQueue.exchange = os.Getenv("EXCHANGE")
	}
	if cfg.messagingQueue.kind == "" {
		cfg.messagingQueue.kind = os.Getenv("EXCHANGE_KIND")
	}
	if cfg.messagingQueue.queue == "" {
		cfg.messagingQueue.queue = os.Getenv("QUEUE")
	}
	if cfg.messagingQueue.routingKey == "" {
		cfg.messagingQueue.routingKey = os.Getenv("ROUTING_KEY")
	}
	if cfg.amqpConnectionURL == "" {
		cfg.amqpConnectionURL = os.Getenv("AMQP_CONNECTION_URL")
	}

	if cfg.amqpConnectionURL == "" {
		cfg.amqpConnectionURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	log.Println("Effective Config:")
	log.Printf("  LLM_ENDPOINT: %s", cfg.llm.endpoint)
	log.Printf("  DB_DSN: %s", cfg.dbDSN)
	log.Printf("  AMQP: %s", cfg.amqpConnectionURL)
	log.Printf("  Queue: exchange=%s kind=%s queue=%s", cfg.messagingQueue.exchange, cfg.messagingQueue.kind, cfg.messagingQueue.queue)

}

func validateConfig(cfg *config) {
	var missingConfigs []string

	if cfg.llm.endpoint == "" {
		missingConfigs = append(missingConfigs, "LLM_ENDPOINT")
	}
	if cfg.llm.authorization == "" {
		missingConfigs = append(missingConfigs, "LLM_AUTHORIZATION_KEY")
	}
	if cfg.dbDSN == "" {
		missingConfigs = append(missingConfigs, "DB_DSN")
	}
	if cfg.messagingQueue.exchange == "" {
		missingConfigs = append(missingConfigs, "EXCHANGE")
	}
	if cfg.messagingQueue.kind == "" {
		missingConfigs = append(missingConfigs, "EXCHANGE_KIND")
	}
	if cfg.messagingQueue.queue == "" {
		missingConfigs = append(missingConfigs, "QUEUE")
	}
	if cfg.amqpConnectionURL == "" {
		missingConfigs = append(missingConfigs, "AMQP_CONNECTION_URL")

	}

	if len(missingConfigs) > 0 {
		log.Fatalf("missing required configs: %v", missingConfigs)
	}
}
