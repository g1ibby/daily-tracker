build:
	mkdir -p functions
	go get ./...
	go build -o functions/updater netlify/updater/updater.go
	go build -o functions/scheduler netlify/scheduler/scheduler.go
