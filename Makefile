EXAMPLE1=1_basic
EXAMPLE2=2_with_kind
EXAMPLE3=3_with_stack
EXAMPLE4=4_custom_error
EXAMPLE5=5_multi_error
EXAMPLE6=6_catch_panic

example%:
	go run -race ./.examples/$(EXAMPLE$(*))/main.go
