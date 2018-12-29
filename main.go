package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cbush06/rpi-golang-test/i2c"
	"github.com/cbush06/rpi-golang-test/i2c/imu"

	"github.com/cbush06/rpi-golang-test/ledarray"
)

var (
	errorImageArgumentMissing = errors.New("[image] mode specified, but no image path found")
)

func main() {
	args := os.Args[1:]

	led := ledarray.NewLed()
	led.Init()
	defer led.Close()

	switch args[0] {
	case "yinyang":
		led.YinYang()
	case "whiteout":
		led.WhiteOut()
	case "image":
		if len(args) < 2 {
			fmt.Println(errorImageArgumentMissing.Error())
			return
		}
		led.FromImage(args[1])
	case "imu":
		i2cConn := i2c.NewI2c()
		imu.AutoDiscover(i2cConn)
	}

	// wait for Ctrl-C, then clear the screen
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Clear
	led.Clear()
}
