package ds18b20

import (
	"fmt"
	"math"
)

// ImpossibleTemperature is returned as a sentinel in error conditions.
const ImpossibleTemperature = Temperature(math.MaxFloat32)

// A Temperature is a temperature value with multiple representations.
type Temperature float32

// Celsius returns a temperature in degrees Celsius.
func (t Temperature) Celsius() float32 {
	return float32(t)
}

// Fahrenheit returns a temperature in degrees Fahrenheit.
func (t Temperature) Fahrenheit() float32 {
	return float32((t * 1.8) + 32.0)
}

func (t Temperature) String() string {
	return fmt.Sprintf("%.3f°C", t)
}
