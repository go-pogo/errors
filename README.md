errors
=======

[![Latest release][latest-release-img]][latest-release-url]
[![Build status][build-status-img]][build-status-url]
[![Go Report Card][report-img]][report-url]
[![Documentation][doc-img]][doc-url]

[latest-release-img]: https://img.shields.io/github/release/go-pogo/errors.svg?label=latest
[latest-release-url]: https://github.com/go-pogo/errors/releases
[build-status-img]: https://github.com/go-pogo/errors/workflows/Go/badge.svg
[build-status-url]: https://github.com/go-pogo/errors/actions?query=workflow%3AGo
[report-img]: https://goreportcard.com/badge/github.com/go-pogo/errors
[report-url]: https://goreportcard.com/report/github.com/go-pogo/errors
[doc-img]: https://godoc.org/github.com/go-pogo/errors?status.svg
[doc-url]: https://pkg.go.dev/github.com/go-pogo/errors


Errors is a Go package designed to make working with (custom) errors an easy task. It supports adding stacktrace frames so tracing the cause of an error is a breeze. 


```sh
go get github.com/go-pogo/errors
```
```go
import "github.com/go-pogo/errors"
```


## Creating an error

## Wrapping existing errors

Error that is created with `errors.New()` and has an additional stack frame captured `errors.Trace()`:
```text
some error: something happened:
    main.doSomething
       .../errors/examples/2_trace/main.go:17
    main.someAction
       .../errors/examples/2_trace/main.go:12
```

A json unmarshalling error that's traced with `errors.Trace()`:
```text
some error: something bad happened while performing someAction:
    main.someAction
        .../errors/examples/3_trace_existing/main.go:22
  - invalid character 'i' looking for beginning of value:
        main.unmarshal
            .../errors/examples/3_trace_existing/main.go:16
```


## Documentation
Additional detailed documentation is available at [go.dev][doc-url]


### Created with
<a href="https://www.jetbrains.com/?from=roeldev" target="_blank"><img src="https://pbs.twimg.com/profile_images/1206615658638856192/eiS7UWLo_400x400.jpg" width="35" /></a>


## License
Copyright © 2019-2020 [Roel Schut](https://roelschut.nl). All rights reserved.

This project is governed by a BSD-style license that can be found in the [LICENSE](LICENSE) file.
