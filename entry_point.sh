#!/bin/sh

migrate -path ./migrations -database "postgres://admin:secret@postgres/mydb?sslmode=disable" up

exec ./bin/llm-log-pipeline