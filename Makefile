build:
	go build main.go
goClean:
	go mod tidy
clean:
	rm *.exe
run:
	go run main.go	
testCoverage:
	go test   -coverpkg=./internal/services/...,./internal/handlers/...,./internal/repositories/...   -coverprofile=coverage.out   ./tests/...
showCoverage:
	go tool cover -html=coverage.out
