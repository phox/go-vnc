package operators

import "testing"

func TestEqualSlicesOfByte(t *testing.T) {
	// Equal.
	if !EqualSlicesOfByte(nil, nil) {
		t.Errorf("EqualSlicesOfByte([]byte): nil == nil")
	}
	if !EqualSlicesOfByte([]byte{}, []byte{}) {
		t.Errorf("EqualSlicesOfByte([]byte): [] == []")
	}
	if !EqualSlicesOfByte([]byte{1, 2, 3}, []byte{1, 2, 3}) {
		t.Errorf("EqualSlicesOfByte([]byte): [1, 2, 3] == [1, 2, 3]")
	}
	// Not equal.
	if EqualSlicesOfByte([]byte{1, 2, 3}, nil) {
		t.Errorf("EqualSlicesOfByte([]byte): [1, 2, 3] != nil")
	}
	if EqualSlicesOfByte([]byte{1, 2, 3}, []byte{}) {
		t.Errorf("EqualSlicesOfByte([]byte): [1, 2, 3] != []")
	}
	if EqualSlicesOfByte([]byte{1, 2, 3}, []byte{4, 5, 6}) {
		t.Errorf("EqualSlicesOfByte([]byte): [1, 2, 3] != [4, 5, 6]")
	}
}
