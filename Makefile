craft:
	codecrafters test
test:
	go test -coverprofile=coverage.out -v ./...
cover:
	go tool cover -html=coverage.out -o coverage.html
	explorer.exe coverage.html