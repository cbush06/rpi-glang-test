package imu

import (
	"errors"
	"fmt"

	"github.com/cbush06/rpi-golang-test/i2c"
	"github.com/cbush06/rpi-golang-test/i2c/imu/lsm9ds1"
)

var (
	errorAutoDiscoveryFailed = errors.New("Auto-discovery of IMU failed")
)

type ImuSettings struct {
	ImuType            uint8
	I2cSlaveAddress    uint8
	I2cBusAddress      uint8
	I2cPressureAddress uint8
	I2cHumidityAddress uint8
}

// AutoDiscover attempts to automatically detect the IMU and configure settings for it.
func AutoDiscover(i *i2c.I2C) (*ImuSettings, error) {
	i.Init()

	result, err := i.ReadRegister(lsm9ds1.LSM9DS1_AG_ADDRESS0, lsm9ds1.AG_WHO_AM_I, 1)
	if err != nil {
		fmt.Printf("Error encountered while attempting to auto discover LSM9DS1 Magnometer at slave address [%Xh], register address [%Xh]: %s\n", lsm9ds1.LSM9DS1_AG_ADDRESS0, lsm9ds1.AG_WHO_AM_I, err.Error())
	} else if len(result) == 1 && result[0] == lsm9ds1.LSM9DS1_AG_ID {
		fmt.Println("LSM9DS1 found!")
		return &ImuSettings{
			ImuType:         ImuTypeLsm9Ds1,
			I2cSlaveAddress: lsm9ds1.LSM9DS1_AG_ADDRESS0,
		}, nil
	}

	return nil, errorAutoDiscoveryFailed
}
