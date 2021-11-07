# run tests with coverage and race detection
test:
	go test -race -cover -coverprofile=coverage.out -v github.com/go-pogo/errors
	go tool cover -func=coverage.out

bench:
	go test -bench=.

example1:
	go run -race ./.examples/1_basic/main.go

example2:
	go run -race ./.examples/2_trace_existing/main.go

example3:
	go run -race ./.examples/3_with_kind/main.go

example4:
	go run -race ./.examples/4_custom_error/main.go

example5:
	go run -race ./.examples/5_multi_error/main.go

example6:
	go run -race ./.examples/6_catch_panic/main.go
