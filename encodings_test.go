package vnc

import (
	"testing"
)

func TestRawEncoding(t *testing.T) {
	e := &RawEncoding{}
	if got, want := e.Type(), Raw; got != want {
		t.Errorf("RawEncoding.Type() got = %v, want = %v", got, want)
	}
}
