# LLM Log Pipeline

**LLM Log Pipeline** is a containerized Go-based service that uses an LLM (e.g., Meta LLaMA 3.1 Instruct, 8B) to automatically analyze and summarize backend logs. It ingests logs through RabbitMQ, processes them with a large language model, and stores insights in a PostgreSQL database for future reference or visualization.

---

##  Features

* Parses raw application logs with stack traces
* Uses LLM to extract:

    * Root cause
    * Severity
    * Time of occurrence
    * Stack trace analysis
    * File and line number
    * Suggested fix
    * Summary
* Stores enriched logs in PostgreSQL
* Runs as a Dockerized microservice
* Easily extendable (e.g., dashboard, alerting)

---

##  Architecture

![Architecture](./docs/llm_log_pipeline.jpg)

---

##  Getting Started

### 1. Clone the Repo

```bash
git clone https://github.com/Daniel-Sogbey/llm_log_pipeline.git
cd llm_log_pipeline
```

### 2. Environment Setup

Add the `.env` file:

Update it with your own values:

```env
DB_DSN=postgres://user:password@host:5432/dbname?sslmode=disable
TOGETHER_API_KEY=your_api_key_here
```

### 3.  Run with Docker

Build and run the service:

```bash
docker build -t llm_log_pipeline .
docker run --env-file .env llm_log_pipeline
```

### 4. (Optional) Docker Compose

Use the provided Docker Compose setup to launch RabbitMQ, PostgreSQL, and the pipeline:

```bash
docker compose up --build
```

* RabbitMQ Management UI: [http://localhost:15672](http://localhost:15672)
* Default credentials: `guest / guest`

---

## How It Works

1. Producer sends raw log entries into RabbitMQ.
2. Exchange fans out the message to queues.
3. Consumer (written in Go) listens to the queue.
4. Consumer calls the LLM API to analyze logs.
5. Stores structured output in PostgreSQL.


### Example Output from LLM:

```json
{
  "analysis": "The log indicates an error occurred when trying to find a record, resulting in a server error.",
  "cause": "The most likely root cause of the error is a missing or incorrect record in the database.",
  "severity": "high",
  "time_of_occurrence": "2025-04-06T17:32:01Z",
  "stacktrace_insight": "The error occurred in the paymentWebHookHandler function in the checkout package.",
  "file": "cmd/api/checkout.go",
  "line_number": "237",
  "summary": "Record not found error",
  "comprehensive_detail": "This error affects the system by preventing the payment webhook handler from functioning correctly, resulting in failed payment processing. From a user's perspective, this may cause payment failures or errors when attempting to process payments. The system will need to handle this error and provide a suitable response to the user.",
  "suggested_way_to_fix": "Check the database for missing or incorrect records, and ensure that the payment webhook handler is correctly configured to handle record not found errors."
}
```


---

## Tech Stack

* Go
* RabbitMQ
* PostgreSQL
* Together.ai (or any LLM with OpenAI-compatible API)
* Docker

---

## CI/CD Pipeline (Coming Soon)

Planned CI/CD capabilities:

* [x] GitHub Actions workflow for automated testing and builds
* [x] Docker image publishing to Docker Hub

Once implemented, this will allow:

* Automated testing on push/PRs
* Docker image builds triggered by changes
* Docker Hub deployments from GitHub

---

##  Extending This Project

You can add:

*  Web Dashboard UI (Next.js, Vue, etc.)
* Alert system for high-severity logs (Slack, Email)
* LLM provider support: OpenAI, OpenRouter, Groq, etc.

---

## Contributing

1. Fork the repo
2. Create your feature branch:

```bash
git checkout -b feature/foo
```

3. Commit changes:

```bash
git commit -am 'Add foo'
```

4. Push to branch:

```bash
git push origin feature/foo
```

5. Submit a pull request

---

##  Credits

* LLM API: Together.ai
* Diagram made with draw\.io
