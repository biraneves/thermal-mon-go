package thermal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateThresholds(t *testing.T) {
	tests := []struct {
		name        string
		warningThr  float64
		criticalThr float64
		wantErr     bool
	}{
		{"valid thresholds", 75.0, 85.0, false},
		{"equal thresholds", 75.0, 75.0, true},
		{"warning greater than critical", 90.0, 85.0, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateThresholds(test.warningThr, test.criticalThr)
			if test.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !test.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}

func TestCheckThresholds(t *testing.T) {
	tests := []struct {
		name        string
		temp        float64
		warningThr  float64
		criticalThr float64
		expected    Status
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CheckThresholds(test.temp, test.warningThr, test.criticalThr)
			if got != test.expected {
				t.Fatalf("expected status %s, got %s", test.expected, got)
			}
		})
	}
}

func TestReadCPUTemperature(t *testing.T) {
	t.Run("valid millidegrees file", func(t *testing.T) {
		dir := t.TempDir()
		file := filepath.Join(dir, "temp_input")
		if err := os.WriteFile(file, []byte("42000\n"), 0o600); err != nil {
			t.Fatalf("write file: %v", err)
		}

		got, err := ReadCPUTemperature(file)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 42.0 {
			t.Fatalf("expected 42.0, got %v", got)
		}
	})

	t.Run("invalid numeric content", func(t *testing.T) {
		dir := t.TempDir()
		file := filepath.Join(dir, "temp_input")
		if err := os.WriteFile(file, []byte("abc"), 0o600); err != nil {
			t.Fatalf("write file: %v", err)
		}

		_, err := ReadCPUTemperature(file)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("missing file", func(t *testing.T) {
		_, err := ReadCPUTemperature("/tmp/this-file-should-not-exist")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}
