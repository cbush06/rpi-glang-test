package i2c

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	// IOCTL Commands re/ I2C
	// For more details, see linux/i2c-dev.h

	ioctlI2cRetries    = 0x0701
	ioctlI2cTimeout    = 0x0702
	ioctlI2cSlave      = 0x0703
	ioctlI2cSlaveForce = 0x0706
	ioctlI2cTenBit     = 0x0704
	ioctlI2cFuncs      = 0x0705
	ioctlI2cRdWr       = 0x0707
	ioctlI2cPec        = 0x0708
	ioctlI2cSmbus      = 0x0720
)

var (
	errorI2CDeviceNotFound       = errors.New("I2C device not found")
	errorI2CDeviceFileNotOpen    = errors.New("I2C device file not open")
	errorI2CReadFailed           = errors.New("I2C read failed")
	errorI2CSlaveSelectFailed    = errors.New("I2C slave selection via IOCTL failed")
	errorI2CRegisterSelectFailed = errors.New("I2C register selection failed")
)

// I2C represents an I2C bus on a linux device
type I2C struct {
	DeviceName            string
	BusID                 uint8
	CurrentSlaveAddress   uint8
	CurrentRegiserAddress uint8
	fileHandle            *os.File
}

// NewI2c returns a new I2C with default settings
func NewI2c() *I2C {
	return &I2C{
		DeviceName:            "bcm2835 I2C adapter",
		BusID:                 1,
		CurrentSlaveAddress:   255,
		CurrentRegiserAddress: 255,
		fileHandle:            nil,
	}
}

// Init attempts to open the system /dev/i2c-* file that corresponds to this I2C device.
func (i *I2C) Init() error {
	if i.fileHandle != nil {
		return nil
	}

	i2cDevice, err := i.getDevice()
	if err != nil {
		fmt.Printf("Could not locate I2C device file for name [%s]\n", i.DeviceName)
		return err
	}

	i.fileHandle, err = os.OpenFile(i2cDevice, os.O_RDWR, 0755)
	if err != nil {
		fmt.Printf("Error encountered opening I2C device file at [%s]: %s\n", i2cDevice, err.Error())
		return err
	}

	return nil
}

// Close releases the /dev/i2c-* file handle assosciated with this device.
func (i *I2C) Close() error {
	if i.fileHandle == nil {
		return nil
	}

	i.CurrentSlaveAddress = 255
	i.CurrentRegiserAddress = 255
	return i.fileHandle.Close()
}

// ReadRegister performs a normal read with register select (for those I2C devices that offer more than one register)
func (i *I2C) ReadRegister(slaveAddr uint8, registerAddr uint8, readLength uint8) ([]byte, error) {
	err := i.selectSlave(slaveAddr)
	if err != nil {
		fmt.Printf("Attempt to select slave [%Xh] failed: %s\n", slaveAddr, err.Error())
		return nil, errorI2CReadFailed
	}

	err = i.selectRegister(registerAddr)
	if err != nil {
		fmt.Printf("Attempt to select register [%Xh] failed: %s\n", registerAddr, err.Error())
		return nil, errorI2CReadFailed
	}

	return i.Read(slaveAddr, readLength)
}

// Read without register select
func (i *I2C) Read(slaveAddr uint8, readLength uint8) ([]byte, error) {
	if i.fileHandle == nil {
		return nil, errorI2CDeviceFileNotOpen
	}

	err := i.selectSlave(slaveAddr)
	if err != nil {
		fmt.Printf("Attempt to select slave [%Xh] failed: %s\n", slaveAddr, err.Error())
		return nil, errorI2CReadFailed
	}

	var data = make([]byte, readLength) // len = 0, cap = readLength

	// Try at least 5 times to read the requested data before giving up
	var total uint8
	var tries uint8
	for total < readLength && tries < 5 {
		bytesRead, err := i.fileHandle.Read(data[total:]) // reposition the slice so data is appended
		if err != nil {
			fmt.Printf("Error encountered reading value from slave [%Xh], failed: %s\n", slaveAddr, err.Error())
			return []byte{}, err
		}
		total += uint8(bytesRead)
		tries++
		time.Sleep(10 * time.Millisecond)
	}

	if total != readLength {
		fmt.Printf("Read from slave [%Xh], register [%Xh] failed\n", slaveAddr, i.CurrentRegiserAddress)
		return []byte{}, errorI2CReadFailed
	}

	return data, nil
}

func (i *I2C) getDevice() (string, error) {
	matches, err := filepath.Glob("/sys/class/i2c-dev/i2c-*")
	if err != nil {
		fmt.Printf("Error listing I2C devices: %s\n", err.Error())
		return "", err
	}

	for _, dir := range matches {
		deviceName, err := ioutil.ReadFile(filepath.Join(dir, "name"))
		if err != nil {
			fmt.Printf("Error reading name of I2C device [%s]: %s\n", dir, err.Error())
			continue
		}

		if i.DeviceName == strings.TrimSpace(string(deviceName)) {
			busID, _ := strconv.Atoi(dir[len(dir)-1:])
			i.BusID = uint8(busID)
			return filepath.Join("/dev", filepath.Base(dir)), nil
		}
	}

	return "", errorI2CDeviceNotFound
}

func (i *I2C) selectSlave(slaveAddr uint8) error {
	if i.CurrentSlaveAddress == slaveAddr {
		return nil
	}
	if i.fileHandle == nil {
		return errorI2CDeviceFileNotOpen
	}

	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, i.fileHandle.Fd(), uintptr(ioctlI2cSlave), uintptr(slaveAddr))
	if err != 0 {
		fmt.Printf("Attempt to select I2C slave [%Xh] via IOCTL failed: %s\n", slaveAddr, err.Error())
		return errorI2CSlaveSelectFailed
	}

	i.CurrentSlaveAddress = slaveAddr
	return nil
}

func (i *I2C) selectRegister(registerAddr uint8) error {
	if registerAddr == i.CurrentRegiserAddress {
		return nil
	}
	if i.fileHandle == nil {
		return errorI2CDeviceFileNotOpen
	}

	if i.CurrentRegiserAddress != registerAddr {
		len, err := i.fileHandle.Write([]byte{registerAddr})
		if err != nil {
			fmt.Printf("Error encountered while writing register address [%Xh]: %s\n", registerAddr, err.Error())
			return errorI2CRegisterSelectFailed
		}
		if len != 1 {
			fmt.Printf("Attempted to write register address [%Xh], but %d byte(s) written instead of 1\n", registerAddr, len)
			return errorI2CRegisterSelectFailed
		}
	}
	i.CurrentRegiserAddress = registerAddr
	return nil
}
