run:
	go run cmd/server/main.go

build:
	cd cmd/server; CGO_ENABLED=0 go build -o ../../bin/go-backup-gdrive

exec:
	./bin/go-backup-gdrive

startapp: build exec