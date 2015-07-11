/*
common.go provides common things that aren't part of the RFB protocol.
*/
package vnc

import "time"

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

// Error returns an VNCError as a string.
func (e VNCError) Error() string {
	return e.s
}

var settleDuration = 25 * time.Millisecond

// Settle returns the UI settle duration.
func Settle() time.Duration {
	return settleDuration
}

// SetSettle changes the UI settle duration.
func SetSettle(s time.Duration) {
	settleDuration = s
}

// settleUI allows the UI to "settle" before the next UI change is made.
func settleUI() {
	time.Sleep(settleDuration)
}
