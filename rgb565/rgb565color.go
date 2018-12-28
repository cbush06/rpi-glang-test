package rgb565

import "image/color"

// Rgb565Color represents an RGB565 value
type Rgb565Color uint16

// New creates a new Color value using the specified component values
func New(red uint8, green uint8, blue uint8) Rgb565Color {
	r := red >> 3
	g := green >> 2
	b := blue >> 3
	return (Rgb565Color(r) << 11) | (Rgb565Color(g) << 5) | Rgb565Color(b)
}

// FromColor converts a color.Color to an Rgb565Color
func FromColor(value color.Color) Rgb565Color {
	r, g, b, _ := value.RGBA()
	return New(uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

// Red returns the red channel value of an Rgb565Color. The lower 3 bits are ignored.
func Red(r Rgb565Color) uint8 {
	return uint8((r & 0xF800) >> 8)
}

// Green returns the green channel value of an Rgb565Color. The lower 3 bits are ignored.
func Green(r Rgb565Color) uint8 {
	return uint8((r & 0x07E0) >> 3)
}

// Blue returns the blue channel value of an Rgb565Color. The lower 3 bits are ignored.
func Blue(r Rgb565Color) uint8 {
	return uint8((r & 0x001F) << 3)
}
