/*
Common provides common things used by multiple source files.
*/
package vnc

const (
	RFBFalse = uint8(iota)
	RFBTrue
)

// VNCError implements error interface.
type VNCError struct {
	s string
}

// NewVNCError returns a custom VNCError error.
func NewVNCError(s string) error {
	return &VNCError{s}
}

func (e VNCError) Error() string {
	return e.s
}
