FROM golang:1.24 AS build-stage

WORKDIR /llm-log-pipeline

COPY go.mod go.sum ./

COPY main.go ./

COPY internal ./internal

COPY migrations ./migrations

COPY entry_point.sh ./

# Build Delve
#RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN #mv /go/bin/dlv /usr/local/bin/

# Install go-migrate CLI
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Move go-migrate to global PATH
RUN mv /go/bin/migrate /usr/local/bin/

RUN GOOS=linux GOARCH=amd64 go build -o ./bin/llm-log-pipeline .

FROM debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /llm-log-pipeline

#COPY --from=build-stage /usr/local/bin/dlv /dev

COPY --from=build-stage /llm-log-pipeline/migrations ./migrations

COPY --from=build-stage /usr/local/bin/migrate /usr/local/bin/

COPY --from=build-stage ./llm-log-pipeline/bin/llm-log-pipeline ./bin/llm-log-pipeline

COPY --from=build-stage ./llm-log-pipeline/entry_point.sh ./

RUN chmod +x ./entry_point.sh

CMD ["./entry_point.sh"]
