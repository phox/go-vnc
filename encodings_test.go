package vnc

// TODO(kward): Fully test the encodings.

import (
	"testing"
)

func TestRawEncoding(t *testing.T) {
	e := &RawEncoding{}
	if got, want := e.Type(), Raw; got != want {
		t.Errorf("incorrect encoding; got = %v, want = %v", got, want)
	}
}

func TestDesktopSizePseudoEncoding(t *testing.T) {
	e := &DesktopSizePseudoEncoding{}
	if got, want := e.Type(), DesktopSizePseudo; got != want {
		t.Errorf("incorrect encoding; got = %v, want = %v", got, want)
	}
}
