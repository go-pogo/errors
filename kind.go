// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

// UnknownKind is the default Kind for errors that are created without a
// distinct Kind.
const UnknownKind Kind = ""

// Kind describes the kind/type of error that has occurred. For example "auth
// error", "unmarshal error", etc. This way errors can be of the same Kind
// but still contain different error messages or additional fields.
// It is recommended to define each Kind as a constant.
type Kind string

// String returns the string representation of Kind.
func (k Kind) String() string { return string(k) }

// Kinder interfaces provide access to a Kind.
type Kinder interface {
	error
	Kind() Kind
}

// WithKind adds Kind to the error.
func WithKind(parent error, kind Kind) Kinder {
	if parent == nil {
		return nil
	}

	switch e := parent.(type) {
	case *kindErr:
		e.kind = kind
		return e

	case Proxy:
		ce := toCommonErr(parent, true)
		ce.kind = kind
		return ce
	}

	return &kindErr{
		error: parent,
		kind:  kind,
	}
}

// GetKind returns the Kind of the error if it implements the Kinder
// interface. If not, it returns UnknownKind.
func GetKind(err error) Kind {
	if e, ok := err.(Kinder); ok {
		return e.Kind()
	}

	return UnknownKind
}

type kindErr struct {
	error
	kind Kind
}

func (e *kindErr) Original() error { return e.error }
func (e *kindErr) Kind() Kind      { return e.kind }
func (e *kindErr) Error() string   { return kindErrMsg(e.error.Error(), e.kind) }

func kindErrMsg(msg string, kind Kind) string {
	if kind == UnknownKind {
		return msg
	}

	return kind.String() + ": " + msg
}
