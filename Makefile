build:
	go build main.go
goClean:
	go mod tidy
clean:
	rm *.exe
run:
	go run main.go	
