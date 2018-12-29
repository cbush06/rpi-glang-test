package imu

import (
	"errors"
)

var (
	errorAutoDiscoveryFailed = errors.New("Auto-discovery of IMU failed")
)

type ImuSettings struct {
	ImuType              uint8
	BusName              string
	BusID                uint8
	I2cAccelGyroAddress  uint8
	I2cMagnometerAddress uint8
	I2cBusAddress        uint8
	I2cPressureAddress   uint8
	I2cHumidityAddress   uint8
}
