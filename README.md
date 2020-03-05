go-errs
=======

[![Latest release][latest-release-img]][latest-release-url]
[![Build status][build-status-img]][build-status-url]
[![Go Report Card][report-img]][report-url]
[![Documentation][doc-img]][doc-url]
![Minimal Go version][go-version-img]

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


## Output examples

Error that is created with `errs.Err()` and wrapped with `errs.Wrap()`:
```text
some error: something happened

Trace:
.../go-errs/examples/2_wrap/main.go:16: main.doSomething():
.../go-errs/examples/2_wrap/main.go:10: main.someAction():
> something happened
```

A json unmarshalling error that's wrapped with `errs.Wrap()`:
```text
invalid character 'i' looking for beginning of value

Trace:
.../go-errs/examples/3_wrap_existing/main.go:18: main.doSomething():
> invalid character 'i' looking for beginning of value
```


## Documentation
Additional detailed documentation is available at [go.dev][doc-url]


### Created with
<a href="https://www.jetbrains.com/?from=roeldev/go-errs" target="_blank"><img src="https://pbs.twimg.com/profile_images/1206615658638856192/eiS7UWLo_400x400.jpg" width="35" /></a>


## License
[GPL-3.0+](LICENSE) © 2019 [Roel Schut](https://roelschut.nl)
