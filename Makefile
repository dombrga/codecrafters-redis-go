craft:
	codecrafters test
test:
	go test -coverprofile=coverage.out -v ./...
# test specific test function
testone:
	go test -v -run ${name} ./...
cover: test
	go tool cover -html=coverage.out -o coverage.html
	explorer.exe coverage.html