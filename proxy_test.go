package errors

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOriginal(t *testing.T) {
	tests := map[string]error{
		"error":     New("original"),
		"std error": stderrors.New("original std error"),
	}

	for name, err := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Same(t, err, Original(Upgrade(err)))
		})
	}
}
