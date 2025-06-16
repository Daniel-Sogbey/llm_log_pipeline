package pubsub

import (
	"database/sql"
	"log"
	"log_processor/internal/data"
	llm2 "log_processor/internal/llm"
)

func (ps *PubSub) ConsumeMessage(exchange, kind, queueName, key string, durable bool, newLLM *llm2.LLM, db *sql.DB) error {

	err := ps.ExchangeDeclare(exchange, kind, durable)
	if err != nil {
		return err
	}

	queue, err := ps.QueueDeclare(queueName, durable)
	if err != nil {
		return err
	}

	err = ps.QueueBind(queue.Name, key, exchange)
	if err != nil {
		return err
	}

	messages, err := ps.channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("recovered in message consumer %v\n", r)
			}
		}()
		for msg := range messages {
			llmLogAnalysis, err := newLLM.AnalyzeLog(llm2.LLMRequestModel{
				Model: newLLM.Model,
				Messages: []llm2.Message{
					{
						Role:    "user",
						Content: llm2.Prompt(string(msg.Body)),
					},
				},
			})

			if err != nil {
				log.Printf("ERROR GETTING LLM LOG ANALYSIS %v\n", err)
			}

			err = data.LogAnalysisDB{DB: db}.ExtractAndInsertLogAnalysis(llmLogAnalysis.Choices[0].Message.Content)
			if err != nil {
				log.Println("Error Extracting and Inserting Log Analysis", err)
			}
		}
	}()

	return nil
}
