package i2c

import (
	"errors"
	"fmt"

	"github.com/cbush06/rpi-golang-test/i2c/imu"
	"github.com/cbush06/rpi-golang-test/i2c/imu/lsm9ds1"
)

var errorLsm9Ds1AutoDiscoveryFailed = errors.New("Unable to find LSM9DS1 via auto discovery")

// LSM9DS1 Accel/Gryo auto discovery
func (i *I2C) autoDiscoverLsm9Ds1() error {
	var result []byte
	var agFound bool
	var magFound bool

	// Try A/G Address 0
	result, _ = i.ReadRegister(lsm9ds1.LSM9DS1_AG_ADDRESS0, lsm9ds1.AG_WHO_AM_I, 1)
	if len(result) == 1 && result[0] == lsm9ds1.LSM9DS1_AG_ID {
		fmt.Printf("AccelGyro found at [%Xh]\n", lsm9ds1.LSM9DS1_AG_ADDRESS0)
		i.Settings.ImuType = imu.ImuTypeLsm9Ds1
		i.Settings.I2cAccelGyroAddress = lsm9ds1.LSM9DS1_AG_ADDRESS0
		agFound = true
	}

	// Try A/G Address 1
	if !agFound {
		result, _ = i.ReadRegister(lsm9ds1.LSM9DS1_AG_ADDRESS1, lsm9ds1.AG_WHO_AM_I, 1)
		if len(result) == 1 && result[0] == lsm9ds1.LSM9DS1_AG_ID {
			fmt.Printf("AccelGyro found at [%Xh]\n", lsm9ds1.LSM9DS1_AG_ADDRESS1)
			i.Settings.ImuType = imu.ImuTypeLsm9Ds1
			i.Settings.I2cAccelGyroAddress = lsm9ds1.LSM9DS1_AG_ADDRESS1
			agFound = true
		}
	}

	// Try Mag Address 0
	fmt.Println("Trying Mag address 0")
	result, _ = i.ReadRegister(lsm9ds1.LSM9DS1_MAG_ADDRESS0, lsm9ds1.MAG_WHO_AM_I, 1)
	if len(result) == 1 && result[0] == lsm9ds1.LSM9DS1_MAG_ID {
		fmt.Printf("Magnetometer found at [%Xh]\n", lsm9ds1.LSM9DS1_MAG_ADDRESS0)
		i.Settings.I2cMagnometerAddress = lsm9ds1.LSM9DS1_MAG_ADDRESS0
		magFound = true
	}

	// Try Mag Address 1
	if !magFound {
		fmt.Println("Trying Mag address 1")
		result, _ = i.ReadRegister(lsm9ds1.LSM9DS1_MAG_ADDRESS1, lsm9ds1.MAG_WHO_AM_I, 1)
		if len(result) == 1 && result[0] == lsm9ds1.LSM9DS1_MAG_ID {
			fmt.Printf("Magnetometer found at [%Xh]\n", lsm9ds1.LSM9DS1_MAG_ADDRESS1)
			i.Settings.I2cMagnometerAddress = lsm9ds1.LSM9DS1_MAG_ADDRESS1
			magFound = true
		}
	}

	if agFound && magFound {
		return nil
	}
	return errorLsm9Ds1AutoDiscoveryFailed
}
