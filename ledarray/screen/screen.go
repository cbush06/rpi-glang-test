package screen

import (
	"encoding/binary"
	"image/color"
	"os"

	"github.com/cbush06/rpi-golang-test/ledarray/rgb565"
	"github.com/cbush06/rpi-golang-test/ledarray/texture"
)

// Screen represents the currently rendered screen on a device (e.g. the Pi Sense-Hat)
type Screen struct {
	texture *texture.Texture
}

// New builds and returns a new, empty screen with the set dimensions
func New(width uint16, height uint16) *Screen {
	return &Screen{
		texture: texture.New(width, height),
	}
}

// NewFrom builds and returns a new screen that uses the provided texture as its backing
func NewFrom(texture *texture.Texture) *Screen {
	return &Screen{
		texture: texture,
	}
}

// At returns an RGBA representation of the pixel at (x, y)
func (s *Screen) At(x uint16, y uint16) color.Color {
	value := s.texture.At(x, y)

	return color.RGBA{
		R: rgb565.Red(value),
		G: rgb565.Green(value),
		B: rgb565.Blue(value),
		A: 0xFF,
	}
}

// Set sets the pixel at (x, y) of the curently applied texture to value
func (s *Screen) Set(x uint16, y uint16, value color.Color) {
	s.texture.Set(x, y, rgb565.FromColor(value))
}

// DrawToDevice writes current texture's pixel values to the device
func (s *Screen) DrawToDevice(d *os.File) {
	d.Seek(0, 0)
	binary.Write(d, binary.LittleEndian, s.texture.GetPixels())
}

// GetTexture returns a pointer to the Texture that backs this Screen
func (s *Screen) GetTexture() *texture.Texture {
	return s.texture
}
