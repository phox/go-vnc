package vnc

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/kward/go-vnc/go/operators"
)

func TestPixelFormatBytes(t *testing.T) {
	tests := []struct {
		pf PixelFormat
		b  []byte
		ok bool
	}{
		//
		// Valid PixelFormats.
		//
		{PixelFormat{BPP: 8, Depth: 8, BigEndian: RFBTrue, TrueColor: RFBFalse},
			[]uint8{8, 8, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, true},
		{PixelFormat{BPP: 8, Depth: 16, BigEndian: RFBTrue, TrueColor: RFBFalse},
			[]uint8{8, 16, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, true},
		{NewPixelFormat(),
			[]uint8{16, 16, 1, 1, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0}, true},
		//
		// Invalid PixelFormats.
		//
		// BPP invalid
		{PixelFormat{BPP: 1, Depth: 1, BigEndian: RFBTrue, TrueColor: RFBFalse},
			[]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false},
		// Depth invalid
		{PixelFormat{BPP: 8, Depth: 1, BigEndian: RFBTrue, TrueColor: RFBFalse},
			[]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false},
		// BPP > Depth
		{PixelFormat{BPP: 16, Depth: 8, BigEndian: RFBTrue, TrueColor: RFBFalse},
			[]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false},
	}

	for _, tt := range tests {
		pf := tt.pf
		b, err := pf.Bytes()
		if err == nil && !tt.ok {
			t.Fatal("PixelFormat.Read() expected error", err)
		}
		if err != nil {
			if verr, ok := err.(*VNCError); !ok {
				t.Errorf("PixelFormat.Read() unexpected %v error: %v", reflect.TypeOf(err), verr)
			}
		}
		if !tt.ok {
			continue
		}
		if got, want := b, tt.b; !operators.EqualSlicesOfByte(got, want) {
			t.Errorf("PixelFormat.Read() got = %v, want = %v", got, want)
		}
	}
}

func TestPixelFormatWrite(t *testing.T) {
	tests := []struct {
		b  []byte
		pf PixelFormat
		ok bool
	}{
		//
		// Valid PixelFormats.
		//
		{[]uint8{8, 8, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			PixelFormat{BPP: 8, Depth: 8, BigEndian: RFBTrue, TrueColor: RFBFalse}, true},
		{[]uint8{8, 16, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			PixelFormat{BPP: 8, Depth: 16, BigEndian: RFBTrue, TrueColor: RFBFalse}, true},
		{[]uint8{32, 32, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			PixelFormat{BPP: 32, Depth: 32, BigEndian: RFBTrue, TrueColor: RFBFalse}, true},
		//
		//
		// Invalid PixelFormats.
		//
	}

	for _, tt := range tests {
		var buf bytes.Buffer
		buf.Write(tt.b)

		pf := PixelFormat{}
		err := pf.Write(&buf)
		if err == nil && !tt.ok {
			t.Fatal("PixelFormat.Write() expected error", err)
		}
		if err != nil {
			if verr, ok := err.(*VNCError); !ok {
				t.Errorf("PixelFormat.Write() unexpected %v error: %v", reflect.TypeOf(err), verr)
			}
		}
		if !tt.ok {
			continue
		}
		if !equalPixelFormat(pf, tt.pf) {
			t.Errorf("PixelFormat.Write() got = %v, want = %v", pf, tt.pf)
		}
	}
}

func equalPixelFormat(got, want PixelFormat) bool {
	gotBytes, err := got.Bytes()
	if err != nil {
		return false
	}
	wantBytes, err := want.Bytes()
	if err != nil {
		return false
	}
	return operators.EqualSlicesOfByte(gotBytes, wantBytes)
}
