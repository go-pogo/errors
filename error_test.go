// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func BenchmarkNew(b *testing.B) {
	disableTraceStack()
	defer enableTraceStack()

	msg := Msg("some err")
	str := "some err"

	b.Run("Msg", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = newCommonErr(msg, false)
		}
	})
	b.Run("Msg ptr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = newCommonErr(&msg, false)
		}
	})
	b.Run("string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = newCommonErr(stderrors.New(str), false)
		}
	})
	b.Run("string to Msg", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = newCommonErr(Msg(str), false)
		}
	})
}

func TestSameErrors(t *testing.T) {
	disableTraceStack()
	defer enableTraceStack()

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
