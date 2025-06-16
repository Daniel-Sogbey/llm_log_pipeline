package main

import (
	"context"
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	"log"
	llm2 "log_processor/internal/llm"
	"log_processor/internal/pubsub"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	llmURL           string
	llmAuthorization string
	dbDSN            string
}

func main() {
	var cfg config
	flag.StringVar(&cfg.llmURL, "llm-url", "https://api.together.xyz/v1/chat/completions", "llm url")
	flag.StringVar(&cfg.llmAuthorization, "llm-auth", os.Getenv("LLM_AUTHORIZATION_KEY"), "llm auth key")
	flag.StringVar(&cfg.dbDSN, "db-dsn", os.Getenv("DB_DSN"), "db dsn")
	flag.Parse()

	log.Println(cfg)

	connectionURL := "amqp://guest:guest@rabbitmq/"
	ps, err := pubsub.NewPubSubConnection(connectionURL)
	if err != nil {
		log.Fatal("Failed to create new pubsub connection. Error: ", err)
	}

	llm := &llm2.LLM{
		URL:           cfg.llmURL,
		Authorization: cfg.llmAuthorization,
	}

	db, err := connectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = ps.ConsumeMessage(pubsub.Exchange, pubsub.Kind, pubsub.Queue, "", pubsub.Durable, llm, db)
	if err != nil {
		log.Fatal("Failed to consume messages. Error: ", err)
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

func connectDB(cfg config) (*sql.DB, error) {
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
