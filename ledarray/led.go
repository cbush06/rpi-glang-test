package ledarray

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cbush06/rpi-golang-test/ledarray/screen"
	"github.com/cbush06/rpi-golang-test/ledarray/texture"
	"golang.org/x/image/bmp"
)

var errorFrameBufferNotFound = errors.New("Frame Buffer Device Not Found")

// Led represents an LED array for an RPi.
type Led struct {
	name        string
	devicePath  string
	frameBuffer *os.File
}

// NewLed creates a new LED and assign it the name "RPi-Sense FB".
func NewLed() *Led {
	return &Led{
		name: "RPi-Sense FB",
	}
}

// Init initializes the LED array by opening its corresponding device file.
func (l *Led) Init() {
	var err error
	l.devicePath, err = getDevice(l.name)
	if err != nil {
		return
	}

	l.frameBuffer, err = openDevice(l.devicePath)
	if err != nil {
		return
	}
}

// Close releases the file handle opened on the LED array's frame buffer device file.
func (l *Led) Close() {
	if l.frameBuffer != nil {
		l.frameBuffer.Close()
	}
}

// Clear deactivates the entire LED array.
func (l *Led) Clear() {
	screen.New(8, 8).DrawToDevice(l.frameBuffer)
}

// YinYang activates half the LED array to white and leaves the other half off.
func (l *Led) YinYang() {
	s := screen.NewFrom(texture.YinYang(8, 8))
	fmt.Printf("Drawing the following to frame buffer: %v\n", s.GetTexture().GetPixels())
	s.DrawToDevice(l.frameBuffer)
}

// WhiteOut activates the entire LED array to white.
func (l *Led) WhiteOut() {
	s := screen.NewFrom(texture.WhiteOut(8, 8))
	fmt.Printf("Drawing the following to frame buffer: %v\n", s.GetTexture().GetPixels())
	s.DrawToDevice(l.frameBuffer)
}

// FromImage colors the LED array based on an 8x8 image specified.
func (l *Led) FromImage(imagePath string) {
	f, err := os.OpenFile(imagePath, os.O_RDONLY, 0444)
	if err != nil {
		fmt.Printf("Error reading from image file [%s]: %s", imagePath, err.Error())
		return
	}

	m, err := bmp.Decode(f)
	if err != nil {
		fmt.Printf("Error decoding image file [%s]: %s", imagePath, err.Error())
		return
	}
	s := screen.NewFrom(texture.FromImage(8, 8, m))
	fmt.Printf("Drawing the following to frame buffer: %v\n", s.GetTexture().GetPixels())
	s.DrawToDevice(l.frameBuffer)
}

func openDevice(name string) (*os.File, error) {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error encountered while opening frame buffer: %s\n", err.Error())
		return nil, err
	}
	fmt.Printf("Frame buffer [%s] opened for writing\n", name)
	return f, nil
}

func getDevice(name string) (string, error) {
	matches, err := filepath.Glob("/sys/class/graphics/fb*")
	if err != nil {
		fmt.Printf("Error listing framebuffer devices: %s\n", err.Error())
		return "", err
	}

	for _, dir := range matches {
		bufferName, err := ioutil.ReadFile(filepath.Join(dir, "name"))
		if err != nil {
			fmt.Printf("Error reading name of framebuffer [%s]: %s\n", dir, err.Error())
			continue
		}

		if name == strings.TrimSpace(string(bufferName)) {
			return filepath.Join("/dev", filepath.Base(dir)), nil
		}
	}

	return "", errorFrameBufferNotFound
}
