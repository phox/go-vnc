package vnc

// TODO(kward): Fully test the encodings.

import (
	"testing"

	"github.com/kward/go-vnc/encodings"
	"github.com/kward/go-vnc/go/operators"
)

func TestEncoding_Marshal(t *testing.T) {
	encs := Encodings{&RawEncoding{}}
	bytes, err := encs.Marshal()
	if err != nil {
		t.Errorf("unexpected error; %s", err)
	}
	if got, want := bytes, []byte{0, 0, 0, 0}; !operators.EqualSlicesOfByte(got, want) {
		t.Errorf("incorrect result; got = %v, want = %v", got, want)
	}
}

func TestRawEncoding_Type(t *testing.T) {
	e := &RawEncoding{}
	if got, want := e.Type(), encodings.Raw; got != want {
		t.Errorf("incorrect encoding; got = %s, want = %s", got, want)
	}
}

func TestRawEncoding_Marshal(t *testing.T) {
	for _, tt := range []struct {
		desc string
		e    *RawEncoding
		data []byte
	}{
		{"empty data",
			&RawEncoding{[]Color{}},
			[]byte{}},
		{"single color",
			&RawEncoding{[]Color{
				Color{&PixelFormat16bit, &ColorMap{}, 0, 127, 7, 0}}},
			[]byte{0, 127}},
		{"multiple colors",
			&RawEncoding{[]Color{
				Color{&PixelFormat16bit, &ColorMap{}, 0, 127, 7, 0},
				Color{&PixelFormat16bit, &ColorMap{}, 0, 32767, 2047, 127}}},
			[]byte{0, 127, 127, 255}},
	} {
		data, err := tt.e.Marshal()
		if err != nil {
			t.Errorf("%s: unexpected error: %s", tt.desc, err)
			continue
		}
		if got, want := data, tt.data; !operators.EqualSlicesOfByte(got, want) {
			t.Errorf("%s: incorrect result; got = %v, want = %v", tt.desc, got, want)
			continue
		}
	}
}

func TestRawEncoding_Read(t *testing.T) {}

func TestDesktopSizePseudoEncoding_Type(t *testing.T) {
	e := &DesktopSizePseudoEncoding{}
	if got, want := e.Type(), encodings.DesktopSizePseudo; got != want {
		t.Errorf("incorrect encoding; got = %s, want = %s", got, want)
	}
}
