# Makefile

build-cli:
	go build -o bin/pkm ./cmd/cli/notes.go

build-server:
	go build -o bin/server ./cmd/server/main.go

install: build-cli
	mv ./bin/pkm ~/go/bin/pkm

run-server:
	./bin/server

conf-prod:
	source ./.env.prod
