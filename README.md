go-errs
=======

[![Latest release][latest-release-img]][latest-release-url]
[![Travis build status][travis-build-img]][travis-build-url]
[![Go Report Card][go-report-img]][go-report-url]
[![GoDoc documentation][go-doc-img]][go-doc-url]

[latest-release-img]: https://img.shields.io/github/release/roeldev/go-errs.svg?label=latest
[latest-release-url]: https://github.com/roeldev/go-errs/releases
[travis-build-img]: https://img.shields.io/travis/roeldev/go-errs.svg
[travis-build-url]: https://travis-ci.org/roeldev/go-errs
[go-report-img]: https://goreportcard.com/badge/github.com/roeldev/go-errs
[go-report-url]: https://goreportcard.com/report/github.com/roeldev/go-errs
[go-doc-img]: https://godoc.org/github.com/roeldev/go-errs?status.svg
[go-doc-url]: https://godoc.org/github.com/roeldev/go-errs

Working with errors in Go can be annoying. This package tries to solve some of these problems by adding information of the origin of the created error.
The package is easily used with custom error types if needed, but should provide plenty of features for you average error message.


## Install
```sh
go get github.com/roeldev/go-errs
```


## Import
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
Additional detailed documentation is available at [godoc.org][go-doc-url]


### Created with
<a href="https://www.jetbrains.com/?from=roeldev/go-errs" target="_blank"><img src="https://pbs.twimg.com/profile_images/809358866442055680/CZjD2GYK_400x400.jpg" width="35" /></a>


## License
[GPL-3.0+](LICENSE) © 2019 [Roel Schut](https://roelschut.nl)
