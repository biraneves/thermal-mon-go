package thermal

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/biraneves/thermal-mon-go/internal/colors"
)

// readCPUTemperature reads the system thermal file and returns degrees Celsius.
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

// checkThresholds evaluates the temperature and logs warnings.
func CheckThresholds(temp, warnThreshold, critThreshold float64) (string, string, error) {
	if warnThreshold >= critThreshold {
		return "", "", errors.New("expected warning threshold to be less than critical threshold")
	}

	currentStatus := "OK"
	colorCode := colors.Green

	if temp >= warnThreshold {
		currentStatus = "WARNING"
		colorCode = colors.Yellow
	}
	if temp >= critThreshold {
		currentStatus = "CRITICAL"
		colorCode = colors.Red
	}

	return currentStatus, colorCode, nil
}
