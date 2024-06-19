// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithTime(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, WithTime(nil, time.Now()))
	})

	for name, wantErr := range provideErrors(true) {
		t.Run(name, func(t *testing.T) {
			wantTime := time.Date(2021, time.December, 28, 0, 1, 2, 3, time.UTC)
			haveErr := WithTime(wantErr, wantTime)
			haveTime, haveHas := GetTime(haveErr)

			assert.Exactly(t, wantTime, haveTime)
			assert.True(t, haveHas)
			assert.ErrorIs(t, haveErr, wantErr)

			t.Run("update", func(t *testing.T) {
				wantTime = time.Now().Add(time.Hour * 24)
				haveErr2 := WithTime(haveErr, wantTime)
				haveTime, haveHas = GetTime(haveErr2)

				assert.Exactly(t, wantTime, haveTime)
				assert.True(t, haveHas)
				assert.Same(t, haveErr, haveErr2)
			})
		})
	}
}
