t: test
test:
	go test -cover -coverprofile=coverage.out -v -race
	go tool cover -func=coverage.out

c: coverage
coverage: test
	go tool cover -html=coverage.out

tidy:
	go mod tidy

examples: example1 example2 example3

example1:
	go run examples/1_basic/main.go

example2:
	go run examples/2_wrap/main.go

example3:
	go run examples/3_wrap_existing/main.go
