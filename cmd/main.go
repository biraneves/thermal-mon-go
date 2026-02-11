package main

import (
	"fmt"
	"os"
	"time"

	"github.com/biraneves/thermal-mon-go/internal/colors"
	"github.com/biraneves/thermal-mon-go/internal/thermal"
)

// Thermal thresholds in Celsius.
const (
	WarningThreshold  = 75.0
	CriticalThreshold = 85.0
	CheckInterval     = 30 * time.Second
	ThermalZonePath   = "/sys/class/hwmon/hwmon1/temp1_input"
)

func main() {
	fmt.Printf("Starting Thermal Monitor at %s\n", time.Now().Format(time.RFC3339))
	fmt.Printf("Thresholds: Warning > %.1f°C, Critical > %.1f°C\n", WarningThreshold, CriticalThreshold)

	ticker := time.NewTicker(CheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		temp, err := thermal.ReadCPUTemperature(ThermalZonePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sError reading temperature: %v%s\n", colors.White, err, colors.Reset)
			continue
		}
		currentStatus, colorCode, err := thermal.CheckThresholds(temp, WarningThreshold, CriticalThreshold)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sInvalid thresholds: %v%s\n", colors.White, err, colors.Reset)
			continue
		}

		fmt.Printf("%s[%s] %s: Current temperature: %.2f°C\n%s",
			colorCode,
			time.Now().Format("2006-01-02 15:04:05"),
			currentStatus,
			temp,
			colors.Reset,
		)
	}
}
