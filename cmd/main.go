package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/biraneves/thermal-mon-go/internal/colors"
	"github.com/biraneves/thermal-mon-go/internal/thermal"
)

// Thermal thresholds in Celsius.
const (
	ThermalZonePath = "/sys/class/hwmon/hwmon1/temp1_input"
)

func main() {
	var (
		warningThreshold  float64
		criticalThreshold float64
		checkingInterval  time.Duration
	)

	flag.Float64Var(&warningThreshold, "w", 75.0, "Warning temperature threshold of operation")
	flag.Float64Var(&criticalThreshold, "c", 85.0, "Critical temperature threshold of operation")
	flag.DurationVar(&checkingInterval, "i", 30*time.Second, "Check interval")
	flag.Parse()

	if len(os.Args) == 1 {
		fmt.Println("You can customize the parameters:")
		fmt.Println("  -w (warning threshold - float)")
		fmt.Println("  -c (critical threshold - float)")
		fmt.Println("  -i (interval - valid duration; e.g.: 20s)")
		fmt.Println("")
	}

	fmt.Printf("Starting Thermal Monitor at %s\n", time.Now().Format(time.RFC3339))
	fmt.Printf("Thresholds: Warning >= %.1f°C, Critical >= %.1f°C\n", warningThreshold, criticalThreshold)
	fmt.Printf("Checking every %v\n\n", checkingInterval)

	ticker := time.NewTicker(checkingInterval)
	defer ticker.Stop()

	for range ticker.C {
		temp, err := thermal.ReadCPUTemperature(ThermalZonePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sError reading temperature: %v%s\n", colors.White, err, colors.Reset)
			continue
		}
		currentStatus, colorCode, err := thermal.CheckThresholds(temp, warningThreshold, criticalThreshold)
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
