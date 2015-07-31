package vnc

// TODO(kward): Fully test the encodings.

import (
	"testing"

	"github.com/kward/go-vnc/go/operators"
)

func TestEncoding_Marshal(t *testing.T) {
	encs := Encodings{&RawEncoding{}}
	bytes, err := encs.Marshal()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got, want := bytes, []byte{0, 0, 0, 0}; !operators.EqualSlicesOfByte(got, want) {
		t.Errorf("incorrect result; got = %v, want = %v", got, want)
	}
}

func TestRawEncoding_Type(t *testing.T) {
	e := &RawEncoding{}
	if got, want := e.Type(), int32(0); got != want {
		t.Errorf("incorrect encoding; got = %v, want = %v", got, want)
	}
}

func TestRawEncoding_Marshal(t *testing.T) {
	tests := []struct {
		e    *RawEncoding
		data []byte
	}{
		{&RawEncoding{[]Color{}},
			[]byte{}},
		{&RawEncoding{[]Color{
			Color{&PixelFormat16bit, &ColorMap{}, 0, 127, 7, 0}}},
			[]byte{0, 127}},
		{&RawEncoding{[]Color{
			Color{&PixelFormat16bit, &ColorMap{}, 0, 127, 7, 0},
			Color{&PixelFormat16bit, &ColorMap{}, 0, 32767, 2047, 127}}},
			[]byte{0, 127, 127, 255}},
	}

	for i, tt := range tests {
		data, err := tt.e.Marshal()
		if err != nil {
			t.Errorf("%v: unexpected error: %v", i, err)
			continue
		}
		if got, want := data, tt.data; !operators.EqualSlicesOfByte(got, want) {
			t.Errorf("%v: incorrect result; got = %v, want = %v", i, got, want)
		}
	}
}

func TestRawEncoding_Read(t *testing.T) {}

func TestDesktopSizePseudoEncoding_Type(t *testing.T) {
	e := &DesktopSizePseudoEncoding{}
	if got, want := e.Type(), int32(-223); got != want {
		t.Errorf("incorrect encoding; got = %v, want = %v", got, want)
	}
}
