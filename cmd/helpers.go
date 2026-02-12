package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/biraneves/thermal-mon-go/internal/colors"
	"github.com/biraneves/thermal-mon-go/internal/thermal"
)

func parseArgs(args []string) (config, error) {
	cfg := config{
		warningThreshold:  defaultWarning,
		criticalThreshold: defaultCritical,
		checkInterval:     defaultInterval,
		thermalZonePath:   defaultThermalZonePath,
	}

	fs := flag.NewFlagSet("thermal-mon", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	fs.Float64Var(&cfg.warningThreshold, "w", cfg.warningThreshold, "Warning temperature threshold in Celsius")
	fs.Float64Var(&cfg.criticalThreshold, "c", cfg.criticalThreshold, "Critical temperature threshold in Celsius")
	fs.DurationVar(&cfg.checkInterval, "i", cfg.checkInterval, "Check interval (example: 20s)")
	fs.StringVar(&cfg.thermalZonePath, "z", cfg.thermalZonePath, "Thermal zone file path")

	if err := fs.Parse(args); err != nil {
		return config{}, err
	}

	if cfg.checkInterval <= 0 {
		return config{}, errors.New("interval must be greater than zero")
	}

	return cfg, nil
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: thermal-mon [-w warning] [-c critical] [-i interval] [-z thermal_zone_path]")
	fmt.Fprintln(w, "  -w float")
	fmt.Fprintln(w, "     warning threshold in Celsius (default 75.0)")
	fmt.Fprintln(w, "  -c float")
	fmt.Fprintln(w, "     critical threshold in Celsius (default 85.0)")
	fmt.Fprintln(w, "  -i duration")
	fmt.Fprintln(w, "     check interval (default 30s, examples: 10s, 1m)")
	fmt.Fprintln(w, "  -z string")
	fmt.Fprintln(w, "     thermal zone file path")
}

func printStartup(cfg config) {
	fmt.Printf("Starting Thermal Monitor at %s\n", time.Now().Format(time.RFC3339))
	fmt.Printf("Thresholds: Warning >= %.1f°C, Critical >= %.1f°C\n", cfg.warningThreshold, cfg.criticalThreshold)
	fmt.Printf("Checking every %v\n", cfg.checkInterval)
	fmt.Printf("Thermal zone: %s\n\n", cfg.thermalZonePath)
}

func runMonitor(cfg config) {
	ticker := time.NewTicker(cfg.checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		temp, err := thermal.ReadCPUTemperature(cfg.thermalZonePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sError reading temperature: %v%s\n", colors.White, err, colors.Reset)
			continue
		}

		status := thermal.CheckThresholds(temp, cfg.warningThreshold, cfg.criticalThreshold)
		colorCode := colorForStatus(status)

		fmt.Printf("%s[%s] %s: Current temperature: %.2f°C%s\n",
			colorCode,
			time.Now().Format("2006-01-02 15:04:05"),
			status,
			temp,
			colors.Reset,
		)
	}
}

func colorForStatus(status thermal.Status) string {
	switch status {
	case thermal.StatusCritical:
		return colors.Red

	case thermal.StatusWarning:
		return colors.Yellow

	default:
		return colors.Green
	}
}
