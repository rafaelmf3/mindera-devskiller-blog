.PHONY: rest-api
rest-api:
	go build -o rest-api 

.PHONY: test
test:
	go test ./...

