errors
======
[![Latest release][latest-release-img]][latest-release-url]
[![Build status][build-status-img]][build-status-url]
[![Go Report Card][report-img]][report-url]
[![Documentation][doc-img]][doc-url]

[latest-release-img]: https://img.shields.io/github/release/go-pogo/errors.svg?label=latest

[latest-release-url]: https://github.com/go-pogo/errors/releases

[build-status-img]: https://github.com/go-pogo/errors/workflows/Test/badge.svg

[build-status-url]: https://github.com/go-pogo/errors/actions/workflows/test.yml

[report-img]: https://goreportcard.com/badge/github.com/go-pogo/errors

[report-url]: https://goreportcard.com/report/github.com/go-pogo/errors

[doc-img]: https://godoc.org/github.com/go-pogo/errors?status.svg

[doc-url]: https://pkg.go.dev/github.com/go-pogo/errors


Package `errors` contains additional functions, interfaces and structs for recording stack frames,
applying basic formatting, working with goroutines, multiple errors and custom error types.

It is inspired by the `golang.org/x/xerrors` package and is designed to be a drop in replacement for
it, as well as the standard library's `errors`
package.

```sh
go get github.com/go-pogo/errors
```

```go
import "github.com/go-pogo/errors"
```

## Stack trace
Every error can track stack trace information. Just wrap it with `errors.WithStack`
and a complete stack trace is captured.

```go
err = errors.WithStack(err)
```

```text
some error: something happened:
    main.doSomething
       .../errors/examples/2_trace/main.go:17
    main.someAction
       .../errors/examples/2_trace/main.go:12
```

## Formatting
Wrap an existing error with `errors.WithFormatter` to upgrade the error to include basic formatting.
Formatting is done using `xerrors.FormatError` and thus the same verbs are supported.

```go
fmt.Printf("%+v", errors.WithFormatter(err))
```

## Catch panics
A convenient function is available to catch panics and store them as an error.

```go
var err error
defer errors.CatchPanic(&err)
```

## Documentation
Additional detailed documentation is available at [pkg.go.dev][doc-url]

## Created with
<a href="https://www.jetbrains.com/?from=go-pogo" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.png" width="35" /></a>

## License
Copyright © 2019-2022 [Roel Schut](https://roelschut.nl). All rights reserved.

This project is governed by a BSD-style license that can be found in the [LICENSE](LICENSE) file.
