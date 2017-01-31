package keys

import (
	"reflect"
	"testing"
)

func TestIntToKeys(t *testing.T) {
	for _, tt := range []struct {
		val  int
		keys Keys
	}{
		{-1234, Keys{Minus, Digit1, Digit2, Digit3, Digit4}},
		{0, Keys{Digit0}},
		{5678, Keys{Digit5, Digit6, Digit7, Digit8}},
	} {
		if got, want := IntToKeys(tt.val), tt.keys; !reflect.DeepEqual(got, want) {
			t.Errorf("IntToKeys(%d) = %v, want %v", tt.val, got, want)
			continue
		}
	}
}
