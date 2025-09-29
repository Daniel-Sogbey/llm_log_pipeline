package pubsub

import (
	"database/sql"
	"log"

	"github.com/Daniel-Sogbey/llm_log_pipeline/internal/data"
	"github.com/Daniel-Sogbey/llm_log_pipeline/internal/llm"
)

func (ps *PubSub) ConsumeMessage(exchange, kind, queueName, key string, durable bool, llmProvider any,
	db *sql.DB) error {

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
			switch llmProvider.(type) {
			case llm.LLM[llm.TogetherAiLLMRequestModel, llm.TogetherAiLLMResponseModel]:
				req := llm.TogetherAiLLMRequestModel{
					Model: llmProvider.(*llm.TogetherAi).Model,
					Messages: []llm.Message{
						{
							Role:    "user",
							Content: llm.Prompt(string(msg.Body)),
						},
					},
				}

				var llmLogAnalysis *llm.TogetherAiLLMResponseModel
				llmLogAnalysis, err = llmProvider.(*llm.TogetherAi).AnalyzeLog(req)

				if err != nil {
					return
				}

				err = data.LogAnalysisDB{DB: db}.ExtractAndInsertLogAnalysis(llmLogAnalysis.Choices[0].Message.Content)
				if err != nil {
					return
				}
			default:
				log.Printf("Invalid llm provider")
			}
		}
	}()

	return nil
}
