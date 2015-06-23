package vnc

import (
	"testing"
)

func TestRawEncoding(t *testing.T) {
	e := NewRawEncoding(nil)
	if got, want := e.Type(), RawEnc; got != want {
		t.Errorf("RawEncoding.Type() got = %v, want = %v", got, want)
	}
}
