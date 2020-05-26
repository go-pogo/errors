package errs

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"golang.org/x/xerrors"
)

type benchmarkFormatErrorHelper struct {
	error
	formatFn func(s fmt.State, v rune)
}

func (h *benchmarkFormatErrorHelper) Format(s fmt.State, v rune) { h.formatFn(s, v) }

func BenchmarkFormatError(b *testing.B) {
	h := benchmarkFormatErrorHelper{
		error: errors.New("error to test benchmark with"),
	}

	b.Run("without pool", func(b *testing.B) {
		h.formatFn = func(s fmt.State, v rune) {
			f := errorFormatter{h.error}
			xerrors.FormatError(f, s, v)
		}

		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = fmt.Sprintf("%+v", h)
		}
	})

	p := sync.Pool{
		New: func() interface{} {
			return errorFormatter{}
		},
	}
	b.Run("with pool", func(b *testing.B) {
		h.formatFn = func(s fmt.State, v rune) {
			f := p.Get().(errorFormatter)
			f.error = h.error

			xerrors.FormatError(f, s, v)
			p.Put(f)
		}

		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = fmt.Sprintf("%+v", h)
		}
	})
}
