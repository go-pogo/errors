 errors
======
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


Package `errors` implements functions to manipulate errors, record stack frames and apply basic formatting to errors. It
is inspired by `golang.org/x/xerrors` and is designed to be a drop in replacement for it, as well as the standard
library's `errors` package. The package contains additional functions, interfaces and structs for working with
goroutines, multiple errors and custom error types.

```sh
go get github.com/go-pogo/errors
```

```go
import "github.com/go-pogo/errors"
```

## Stack trace

Every error can track stack trace information. Just wrap it with `errors.Trace()` and an addition stack frame is
captured and stored within the error.

```text
some error: something happened:
    main.doSomething
       .../errors/examples/2_trace/main.go:17
    main.someAction
       .../errors/examples/2_trace/main.go:12
```

## Documentation

Additional detailed documentation is available at [pkg.go.dev][doc-url]

## Created with

<a href="https://www.jetbrains.com/?from=go-pogo" target="_blank"><img src="https://pbs.twimg.com/profile_images/1206615658638856192/eiS7UWLo_400x400.jpg" width="35" /></a>

## License

Copyright © 2019-2021 [Roel Schut](https://roelschut.nl). All rights reserved.

This project is governed by a BSD-style license that can be found in the [LICENSE](LICENSE) file.
