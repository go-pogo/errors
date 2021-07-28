// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"

	"github.com/go-pogo/errors/internal"
)

func TestSameErrors(t *testing.T) {
	internal.DisableCaptureFrames()
	defer internal.EnableCaptureFrames()

	cause := xerrors.New("cause of error")
	tests := map[string]map[string][2]error{
		"New&Newf": {
			"empty": {New(""), Newf("")},
			"message only": {
				New("some `foo` happened"),
				Newf("some `%s` happened", "foo"),
			},
		},
		"Wrap&Wrapf": {
			"empty": {
				Wrap(cause, ""),
				Wrapf(cause, ""),
			},
			"message only": {
				Wrap(cause, "some `foo` happened"),
				Wrapf(cause, "some `%s` happened", "foo"),
			},
		},
	}

	for group, ts := range tests {
		t.Run(group, func(t *testing.T) {
			for name, errs := range ts {
				t.Run(name, func(t *testing.T) {
					assert.Equal(t, errs[0].Error(), errs[1].Error())
				})
			}
		})
	}
}

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

func TestUpgrade(t *testing.T) {
	internal.DisableCaptureFrames()
	defer internal.EnableCaptureFrames()

	msg := "a really important err msg"
	kind := Kind("some kind")

	tests := map[string]struct {
		err error
		fn  func(want *commonErr, err error)
	}{
		"common error": {
			err: New(msg),
			fn: func(want *commonErr, err error) {
				want.error = stderrors.New(msg)
			},
		},
		"std error": {
			err: stderrors.New(msg),
			fn: func(want *commonErr, err error) {
				want.error = err
				want.upgrade = true
			},
		},
		"common error with kind": {
			err: WithKind(New(msg), kind),
			fn: func(want *commonErr, err error) {
				want.error = stderrors.New(msg)
				want.kind = kind
			},
		},
		"std error with kind": {
			err: WithKind(stderrors.New(msg), kind),
			fn: func(want *commonErr, err error) {
				want.error = stderrors.New(msg)
				want.upgrade = true
				want.kind = kind
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			want := toCommonErr(nil, false)
			tc.fn(want, tc.err)

			assert.Exactly(t, want, Upgrade(tc.err))
		})
	}
}

func TestCommonErr_GoString(t *testing.T) {
	msg := "just some error message"
	assert.Equal(t,
		fmt.Sprintf("&\"%s\".commonErr{error:%#v}", fullPkgName, stderrors.New(msg)),
		fmt.Sprintf("%#v", New(msg)),
	)
}
