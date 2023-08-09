
build-cli:
	go build -o bin/cmd ./cmd/cli/notes.go

build-server:
	go build -o bin/server ./cmd/server/main.go

install:
	go install ./cmd/cli/notes.go

run-server:
	./bin/server

