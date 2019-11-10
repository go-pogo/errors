t: test
test:
	go test -cover -coverprofile=coverage.out -v -race
	go tool cover -func=coverage.out

c: coverage
coverage: test
	go tool cover -html=coverage.out

tidy:
	go mod tidy
