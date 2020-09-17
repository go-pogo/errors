package errors

import (
	stderrors "errors"
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

type benchmarkFormatErrorHelper struct {
	error
	formatFn func(s fmt.State, v rune)
}

func (h *benchmarkFormatErrorHelper) Format(s fmt.State, v rune) { h.formatFn(s, v) }

func BenchmarkFormatError(b *testing.B) {
	h := benchmarkFormatErrorHelper{
		error: stderrors.New("error to test benchmark with"),
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

func TestFormatError(t *testing.T) {
	tests := map[string]struct {
		setup      func() error
		traceLines []int
	}{
		"error": {
			setup: func() error {
				return New(UnknownKind, "some err")
			},
			traceLines: []int{66},
		},
		"primitive": {
			setup: func() error {
				return Trace(stderrors.New("primitive"))
			},
			traceLines: []int{72},
		},
		"traced error": {
			setup: func() error {
				err := New(UnknownKind, "another err")
				return Trace(err)
			},
			traceLines: []int{78, 79},
		},
		"multi error": {
			setup: func() error {
				err1 := New(UnknownKind, "err1")
				err2 := New(UnknownKind, "err2")
				return Combine(err1, err2)
			},
			traceLines: []int{85, 86, 87},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.setup()
			str := fmt.Sprintf("%+v", err)

			for _, line := range tc.traceLines {
				assert.Contains(t, str, "format_test.go:"+strconv.Itoa(line))
			}
		})
	}
}
