package thermal

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Status string

const (
	StatusOK       Status = "OK"
	StatusWarning  Status = "WARNING"
	StatusCritical Status = "CRITICAL"
)

// ReadCPUTemperature reads the system thermal file and returns degrees Celsius.
func ReadCPUTemperature(filePath string) (float64, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read thermal file: %w", err)
	}

	// Value is usually in millidegrees Celsius.
	raw := strings.TrimSpace(string(data))
	milliDeg, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid temperature format: %w", err)
	}

	return milliDeg / 1000.0, nil
}

// ValidateThresholds validates temperature thresholds once at startup.
func ValidateThresholds(warnThreshold, critThreshold float64) error {
	if warnThreshold >= critThreshold {
		return errors.New("warning threshold must be lower than critical threshold")
	}
	return nil
}

// CheckThresholds evaluates the current status based on thresholds.
func CheckThresholds(temp, warnThreshold, critThreshold float64) Status {
	if temp >= critThreshold {
		return StatusCritical
	}
	if temp >= warnThreshold {
		return StatusWarning
	}
	return StatusOK
}
