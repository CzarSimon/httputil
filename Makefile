test:
	go vet ./...
	go test ./...
	gosec ./...