// Package ds18b20 implements routines for working with DS18B20 temperature
// probes on a Raspberry Pi.
package ds18b20

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Operating system resource names.
const (
	DevicesPath     = "/sys/bus/w1/devices/"
	Modprobe        = "/sbin/modprobe"
	ModW1Therm      = "w1-therm"
	ModW1GPIO       = "w1-gpio"
	MasterBusPrefix = "w1_bus_master"
	SensorPrefix    = "28-"
)

// Sensing errors.
var (
	ErrClosed        = errors.New("probe file descriptor is closed")
	ErrNoBus         = errors.New("1-Wire master bus not present")
	ErrNoSlaves      = errors.New("no temperature probes found")
	ErrNotFound      = errors.New("specified temperature probe not found")
	ErrCRC           = errors.New("CRC error")
	ErrInvalidOutput = errors.New("could not parse sensor output")
)

// ID is a DS18B20 sensor identifier.
type ID string

// Ensure loads the w1-gpio and w1-therm modules are loaded and checks that the
// 1-Wire master bus is ready.
func Ensure() error {
	// Load modules.
	if err := exec.Command(Modprobe, ModW1GPIO).Run(); err != nil {
		return err
	}
	if err := exec.Command(Modprobe, ModW1Therm).Run(); err != nil {
		return err
	}

	// Check for master bus device file.
	devices, err := ioutil.ReadDir(DevicesPath)
	if err != nil {
		return err
	}
	for _, device := range devices {
		if strings.HasPrefix(device.Name(), MasterBusPrefix) {
			return nil
		}
	}
	return ErrNoBus
}

// Sensors returns a listing of available sensor IDs.
func Sensors() ([]ID, error) {
	files, err := ioutil.ReadDir(DevicesPath)
	if err != nil {
		return nil, err
	}

	var sensors []ID
	for _, file := range files {
		if strings.HasPrefix(file.Name(), SensorPrefix) {
			if (file.Mode() & os.ModeSymlink) == os.ModeSymlink {
				sensors = append(sensors, ID(file.Name()))
			}
		}
	}
	return sensors, nil
}
