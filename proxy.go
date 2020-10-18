package errors

// A Proxy adds additional context to an existing error, it is not an error by
// itself. It provides a method to retrieve the original error so it can be
// properly matched using Is and As.
type Proxy interface {
	// Original returns the Original error that resides in the Proxy.
	Original() (original error)
}

// Original returns the Original error if err is a Proxy. otherwise it will
// return the given error err.
func Original(err error) error {
	p, ok := err.(Proxy)
	if !ok {
		return err
	}

	return p.Original()
}
