package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/image/bmp"

	"github.com/cbush06/rpi-golang-test/screen"
	"github.com/cbush06/rpi-golang-test/texture"
)

var (
	errorFrameBufferNotFound  = errors.New("Frame Buffer Device Not Found")
	errorImageArgumentMissing = errors.New("[image] mode specifid, but no image path found")
)

func main() {
	args := os.Args[1:]

	devPath, err := getDevice("RPi-Sense FB")
	if err != nil {
		return
	}

	fb, err := openDevice(devPath)
	if err != nil {
		return
	}

	// f, _ := os.OpenFile("/home/cbush/bytes", os.O_WRONLY|os.O_CREATE, 0644)
	// err = binary.Write(f, binary.LittleEndian, s.GetTexture().GetPixels())
	// if err != nil {
	// 	fmt.Printf("Error enountered while writing to file: %s", err.Error())
	// }
	// f.Close()

	defer fb.Close()

	switch args[0] {
	case "yinyang":
		yinYang(fb)
	case "whiteout":
		whiteOut(fb)
	case "image":
		if len(args) < 2 {
			fmt.Println(errorImageArgumentMissing.Error())
			return
		}
		fromImage(args[1], fb)
	}

	// wait for Ctrl-C, then clear the screen
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Clear
	screen.New(8, 8).DrawToDevice(fb)
}

func yinYang(fb *os.File) {
	s := screen.NewFrom(texture.YinYang(8, 8))
	fmt.Printf("Drawing the following to frame buffer: %v\n", s.GetTexture().GetPixels())
	s.DrawToDevice(fb)
}

func whiteOut(fb *os.File) {
	s := screen.NewFrom(texture.WhiteOut(8, 8))
	fmt.Printf("Drawing the following to frame buffer: %v\n", s.GetTexture().GetPixels())
	s.DrawToDevice(fb)
}

func fromImage(imagePath string, fb *os.File) {
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
	s.DrawToDevice(fb)
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
		fmt.Println(dir)
		b, err := ioutil.ReadFile(filepath.Join(dir, "name"))
		if err != nil {
			fmt.Printf("Error reading name of framebuffer [%s]: %s\n", dir, err.Error())
			continue
		}

		fmt.Printf("Frame Buffer Name of [%s]: %s\n", dir, b)

		fbName := strings.TrimSpace(string(b))
		if name == fbName {
			return filepath.Join("/dev", filepath.Base(dir)), nil
		}
	}

	return "", errorFrameBufferNotFound
}
