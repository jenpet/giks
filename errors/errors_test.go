package errors

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsExitErrors(t *testing.T) {
	exitErrorTests := []struct {
		name           string
		inErr          error
		expectedResult bool
	}{
		{
			"random error",
			errors.New("test"),
			false,
		},
		{
			"warning error",
			NewWarningErrorf("foo %s", "bar"),
			true,
		},
		{
			"error nil",
			nil,
			false,
		},
	}
	for _, tt := range exitErrorTests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedResult, IsWarningError(tt.inErr), "expected result does not match input")
		})
	}
}
