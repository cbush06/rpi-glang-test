package texture

import (
	"errors"
	"fmt"
	"image"

	"github.com/cbush06/rpi-golang-test/ledarray/rgb565"
)

var errorsIncompatibleSize = errors.New("The size of the provided image was not at least 8x8")

// Texture represents a matrix of color values that can be applied to a frame buffer
type Texture struct {
	width  uint16
	height uint16
	pixels []rgb565.Rgb565Color
}

// New creates a new Texture and returns the reference to it
func New(width uint16, height uint16) *Texture {
	t := &Texture{
		width:  width,
		height: height,
		pixels: make([]rgb565.Rgb565Color, width*height),
	}

	for i := range t.pixels {
		t.pixels[i] = rgb565.New(0x00, 0x00, 0x00)
	}

	return t
}

// WhiteOut creates a Texture with all its pixels set to white.
func WhiteOut(width uint16, height uint16) *Texture {
	t := New(width, height)
	for i := range t.pixels {
		t.pixels[i] = rgb565.New(0xFF, 0xFF, 0xFF)
	}
	return t
}

// YinYang creates a Texture with its top half of pixels white and lower half black.
func YinYang(width uint16, height uint16) *Texture {
	t := New(width, height)
	for y := uint16(0); y < height/2; y++ {
		for x := uint16(0); x < width; x++ {
			t.Set(x, y, rgb565.New(0xFF, 0xFF, 0xFF))
		}
	}
	for y := uint16(height/2) + 1; y < height; y++ {
		for x := uint16(0); x < width; x++ {
			t.Set(x, y, rgb565.New(0x00, 0x00, 0x00))
		}
	}
	return t
}

// FromImage loads the current texture from the provided image
func FromImage(width uint16, height uint16, m image.Image) *Texture {
	b := m.Bounds()

	// m must be at least 8x8
	if b.Max.X < 8 || b.Max.Y < 8 {
		fmt.Println(errorsIncompatibleSize.Error())
		return New(width, height)
	}

	// set pixels in current texture to the upper left 8x8 pixels of the provided image
	t := New(width, height)
	var x, y uint16
	for y = 0; y < 8; y++ {
		for x = 0; x < 8; x++ {
			p := m.At(int(x), int(y))
			t.Set(x, y, rgb565.FromColor(p))
		}
	}

	return t
}

// Set sets the value of the pixel at locatin (x, y) with the specified value
func (t *Texture) Set(x uint16, y uint16, value rgb565.Rgb565Color) {
	t.pixels[y*t.width+x] = value
}

// At returns the value ofthe pixel at location (x, y)
func (t *Texture) At(x uint16, y uint16) rgb565.Rgb565Color {
	return t.pixels[y*t.width+x]
}

// GetWidth return the widths of this texture
func (t *Texture) GetWidth() uint16 {
	return t.width
}

// GetHeight returns the height of this texture
func (t *Texture) GetHeight() uint16 {
	return t.height
}

// GetPixels returns the array of uint16 RGB565 formatted pixel values
func (t *Texture) GetPixels() []rgb565.Rgb565Color {
	newArray := make([]rgb565.Rgb565Color, len(t.pixels))
	copy(newArray, t.pixels)
	return newArray
}
