package internal

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHandleError(t *testing.T) {
	type args struct {
		name string
		err  error
	}

	test := &args{
		"it does nothing if nil", nil,
	}
	t.Run(test.name, func(t *testing.T) {
		require.NotPanics(t, func() { HandleError(test.err) })
	})

	test = &args{
		"it panics if error passed in", errors.New("an error"),
	}
	t.Run(test.name, func(t *testing.T) {
		require.Panics(t, func() { HandleError(test.err) })
	})
}
