package ds18b20

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Verify interfaces.
var (
	_ io.Closer = &Probe{}
)

// A Probe represents a single DS18B20 sensor attached to the W1 master bus.
type Probe struct {
	ID ID

	fd *os.File
}

// New constructs a new probe by opening the corresponding device file.
func New(id ID) (*Probe, error) {
	fd, err := os.Open(filepath.Join(filepath.Join(DevicesPath, string(id)), "w1_slave"))
	if err != nil {
		return nil, err
	}
	return &Probe{
		ID: id,
		fd: fd,
	}, nil
}

// Sense reads the probe's temperature.
func (p *Probe) Sense() (Temperature, error) {
	// Re-seek to the beginning of the file to signal the hardware device to send
	// a new reading.
	if _, err := p.fd.Seek(0, 0); err != nil {
		return ImpossibleTemperature, err
	}
	reading, err := ioutil.ReadAll(p.fd)
	if err != nil {
		return ImpossibleTemperature, err
	}

	lines := strings.Split(strings.TrimSpace(string(reading)), "\n")
	if len(lines) != 2 {
		return ImpossibleTemperature, ErrInvalidOutput
	}

	// Check CRC.
	if !strings.HasSuffix(lines[0], "YES") {
		return ImpossibleTemperature, ErrCRC
	}

	// Read temperature.
	idx := strings.Index(lines[1], "t=")
	if idx == -1 {
		return ImpossibleTemperature, ErrInvalidOutput
	}
	temp, err := strconv.Atoi(lines[1][idx+2:])
	if err != nil {
		return ImpossibleTemperature, err
	}

	return Temperature(float32(temp) / 1000.0), nil
}

// Close the underlying device file for this probe.
func (p *Probe) Close() error {
	return p.fd.Close()
}
