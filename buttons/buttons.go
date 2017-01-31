package buttons

// Button represents a mask of pointer presses/releases.
type Button uint8

// All available button mask components.
const (
	Left Button = 1 << iota
	Middle
	Right
	Four
	Five
	Six
	Seven
	Eight
	None Button = 0
)

func Mask(button Button) uint8 {
	return uint8(button)
}
