go-errs
=======

[![Latest release][latest-release-img]][latest-release-url]
[![Build status][build-status-img]][build-status-url]
[![Go Report Card][report-img]][report-url]
[![Documentation][doc-img]][doc-url]

[latest-release-img]: https://img.shields.io/github/release/roeldev/go-errs.svg?label=latest
[latest-release-url]: https://github.com/roeldev/go-errs/releases
[build-status-img]: https://github.com/roeldev/go-errs/workflows/Go/badge.svg
[build-status-url]: https://github.com/roeldev/go-errs/actions?query=workflow%3AGo
[report-img]: https://goreportcard.com/badge/github.com/roeldev/go-errs
[report-url]: https://goreportcard.com/report/github.com/roeldev/go-errs
[doc-img]: https://godoc.org/github.com/roeldev/go-errs?status.svg
[doc-url]: https://pkg.go.dev/github.com/roeldev/go-errs


Errs is a Go package designed to make working with (custom) errors an easy task. It supports adding stacktrace frames so tracing the cause of an error is a breeze. 


```sh
go get github.com/roeldev/go-errs
```
```go
import "github.com/roeldev/go-errs"
```


## Creating an error

## Wrapping existing errors

Error that is created with `errs.New()` and has an additional stack frame captured `errs.Trace()`:
```text
some error: something happened:
    main.doSomething
       .../go-errs/examples/2_trace/main.go:17
    main.someAction
       .../go-errs/examples/2_trace/main.go:12
```

A json unmarshalling error that's traced with `errs.Trace()`:
```text
some error: something bad happened while performing someAction:
    main.someAction
        .../go-errs/examples/3_trace_existing/main.go:22
  - invalid character 'i' looking for beginning of value:
        main.unmarshal
            .../go-errs/examples/3_trace_existing/main.go:16
```


## Documentation
Additional detailed documentation is available at [go.dev][doc-url]


### Created with
<a href="https://www.jetbrains.com/?from=roeldev/go-errs" target="_blank"><img src="https://pbs.twimg.com/profile_images/1206615658638856192/eiS7UWLo_400x400.jpg" width="35" /></a>


## License
[GPL-3.0+](LICENSE) © 2019-2020 [Roel Schut](https://roelschut.nl)
