run:
	go run main.go

build:
	cd cmd/server; go build -o ../../bin/go-backup-gdrive

exec:
	./bin/go-backup-gdrive

startapp: build exec