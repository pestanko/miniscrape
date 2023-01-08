package collut

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptGetDefault(t *testing.T) {
	nonNil := "non-nil"
	tests := []struct {
		name  string
		value *string
		defl  string
		exp   string
	}{
		{
			name:  "empty",
			value: nil,
			defl:  "",
			exp:   "",
		},
		{
			name:  "non-nil value",
			value: &nonNil,
			defl:  "default",
			exp:   nonNil,
		},
		{
			name:  "nil value, should return default",
			value: nil,
			defl:  "default",
			exp:   "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.exp, PtrGetDefault(tt.value, tt.defl))
		})
	}
}
