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

// KindGetter interfaces provide access to a Kind.
type KindGetter interface {
	error
	Kind() Kind
}

// WithKind adds Kind to the error.
func WithKind(parent error, kind Kind) KindGetter {
	if parent == nil {
		return nil
	}

	if e, ok := parent.(kindGetterSetter); ok {
		e.setKind(kind)
		return e
	}
	if _, ok := parent.(OriginalGetter); ok {
		ce := upgrade(parent)
		ce.setKind(kind)
		return ce
	}

	return &kindErr{
		error: parent,
		kind:  kind,
	}
}

// GetKind returns the Kind of the error if it implements the KindGetter
// interface. If not, it returns UnknownKind.
func GetKind(err error) Kind { return GetKindOr(err, UnknownKind) }

// GetKindOr returns the Kind of the error if it implements the KindGetter
// interface. If not, it returns the provided value or.
func GetKindOr(err error, or Kind) Kind {
	if e, ok := err.(KindGetter); ok {
		return e.Kind()
	}
	return or
}

type kindGetterSetter interface {
	KindGetter
	setKind(k Kind)
}

type kindErr struct {
	error
	kind Kind
}

func (ce *commonErr) setKind(k Kind) { ce.kind = k }
func (e *kindErr) setKind(k Kind)    { e.kind = k }

func (e *kindErr) Original() error { return e.error }
func (e *kindErr) Kind() Kind      { return e.kind }
func (e *kindErr) Error() string   { return errMsg(e.error.Error(), e.kind, 0) }
