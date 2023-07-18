
build-cli:
	go build -o bin/cmd ./cmd/cli/notes.go

build-server:
	go build -o bin/server ./cmd/server/main.go

run-server:
	./bin/server

