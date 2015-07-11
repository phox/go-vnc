/*
The operators package provides additional operators that aren't part of the
base Go installation.
*/

package operators

// EqualSlicesOfByte compares two byte slices for equality.
func EqualSlicesOfByte(x, y []byte) bool {
	// Special cases.
	switch {
	case len(x) != len(y):
		return false
	}
	for i, v := range x {
		if v != y[i] {
			return false
		}
	}
	return true
}
