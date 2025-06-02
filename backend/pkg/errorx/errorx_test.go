package errorx_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/moLIart/gomoku-backend/pkg/errorx"
)

func TestWrap_NilError(t *testing.T) {
	err := errorx.Wrap(nil, "some message")
	assert.Nil(t, err)
}

func TestWrap_WithError(t *testing.T) {
	origErr := errors.New("original error")
	wrapped := errorx.Wrap(origErr, "context message")
	require.Error(t, wrapped)
	assert.Contains(t, wrapped.Error(), "context message")
	assert.Contains(t, wrapped.Error(), "original error")
	assert.True(t, errors.Is(wrapped, origErr))
}

func TestMustNoError_NoError(t *testing.T) {
	assert.NotPanics(t, func() {
		errorx.MustNoError(nil, "should not panic")
	})
}
