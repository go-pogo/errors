package errors

// UnknownKind is used for errors that are created without a distinct `Kind`.
const UnknownKind Kind = ""

// ErrorWithKind interfaces provide access to a `Kind`.
type ErrorWithKind interface {
	error
	Kind() Kind
}

// Kind describes the kind/type of error that has occurred. For example "auth
// error", "unmarshal error", etc. This way errors can be of the same `Kind`
// but still contain different error messages or additional fields.
// It is recommended to define each `Kind` as a constant.
type Kind string

// String returns the string representation of `Kind`.
func (k Kind) String() string { return string(k) }

// GetKind returns the `Kind` of the error if it implements the `ErrorWithKind`
// interface. If not, it returns `UnknownKind`.
func GetKind(err error) Kind {
	if e, ok := err.(ErrorWithKind); ok {
		return e.Kind()
	}
	return UnknownKind
}
