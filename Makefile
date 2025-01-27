.PHONY: docs
docs:
	@swag init -g ./cmd/main.go && swag fmt

.PHONY: build
build:
	@go build -o ./bin/main ./cmd/

.PHONY: dev
dev:
	@go run ./cmd/

.PHONY: run
run:
	@./bin/main

.PHONY: devbuildsend
devbuildsend:
	@go build -o ./bin/main ./cmd/
	@ssh -t juallagi@158.140.163.172 -p 2222 sudo systemctl stop filestore
	@scp -P 2222 ./bin/main juallagi@158.140.163.172:filestore/filestore
	@ssh -t juallagi@158.140.163.172 -p 2222 sudo systemctl start filestore

.PHONY: prodbuildsend
prodbuildsend:
	@go build -o ./bin/main ./cmd/
	@ssh delogic@103.183.75.231 sudo systemctl stop automate
	@scp ./bin/main delogic@103.183.75.231:automate/
	@ssh delogic@103.183.75.231 sudo systemctl start automate