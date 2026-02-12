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
	defaultThermalZonePath = "/sys/class/hwmon/hwmon1/temp1_input"
	defaultWarning         = 75.0
	defaultCritical        = 85.0
	defaultInterval        = 30 * time.Second
)

type config struct {
	warningThreshold  float64
	criticalThreshold float64
	checkInterval     time.Duration
	thermalZonePath   string
}

func main() {
	cfg, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sInvalid arguments: %v%s\n", colors.White, err, colors.Reset)
		printUsage(os.Stderr)
		os.Exit(2)
	}

	if err := thermal.ValidateThresholds(cfg.warningThreshold, cfg.criticalThreshold); err != nil {
		fmt.Fprintf(os.Stderr, "%sInvalid thresholds: %v%s\n", colors.White, err, colors.Reset)
		os.Exit(2)
	}

	printStartup(cfg)
	runMonitor(cfg)
}
