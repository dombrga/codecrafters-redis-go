craft:
	codecrafters test
test:
	go test -v -coverprofile=coverage.out -coverpkg=./... ./...
# test specific test function
testone:
	go test -v -run ${name} ./...
cover: test
	go tool cover -html=coverage.out -o coverage.html
	explorer.exe coverage.html