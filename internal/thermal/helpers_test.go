package thermal

import (
	"strings"
	"testing"

	"github.com/biraneves/thermal-mon-go/internal/colors"
)

func TestCheckThresholds(t *testing.T) {
	tests := []struct {
		name        string
		temp        float64
		warningThr  float64
		criticalThr float64
		expStatus   string
		expColor    string
		expError    bool
	}{
		{
			name:        "temperature under warning threshold",
			temp:        40.0,
			warningThr:  75.0,
			criticalThr: 85.0,
			expStatus:   "OK",
			expColor:    sanitizeColorCode(colors.Green),
		},
		{
			name:        "temperature is warning threshold",
			temp:        75.0,
			warningThr:  75.0,
			criticalThr: 85.0,
			expStatus:   "WARNING",
			expColor:    sanitizeColorCode(colors.Yellow),
		},
		{
			name:        "temperature between warning and critical threshold",
			temp:        80.0,
			warningThr:  75.0,
			criticalThr: 85.0,
			expStatus:   "WARNING",
			expColor:    sanitizeColorCode(colors.Yellow),
		},
		{
			name:        "temperature is critical threshold",
			temp:        85.0,
			warningThr:  75.0,
			criticalThr: 85.0,
			expStatus:   "CRITICAL",
			expColor:    sanitizeColorCode(colors.Red),
		},
		{
			name:        "temperature greater than critical threshold",
			temp:        90.0,
			warningThr:  75.0,
			criticalThr: 85.0,
			expStatus:   "CRITICAL",
			expColor:    sanitizeColorCode(colors.Red),
		},
		{
			name:        "warning temperature equal to critical temperature",
			temp:        40.0,
			warningThr:  75.0,
			criticalThr: 75.0,
			expStatus:   "",
			expColor:    "",
			expError:    true,
		},
		{
			name:        "warning temperature greater than critical temperature",
			temp:        40.0,
			warningThr:  90.0,
			criticalThr: 85.0,
			expStatus:   "",
			expColor:    "",
			expError:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotStatus, gotColor, gotErr := CheckThresholds(test.temp, test.warningThr, test.criticalThr)
			if test.expError {
				if gotErr == nil {
					t.Fatalf("Error expected, got nil")
				}
				if gotStatus != "" {
					t.Fatalf("Expected empty status")
				}
				if gotColor != "" {
					t.Fatalf("Expected empty color code")
				}
			} else {
				if gotStatus != test.expStatus {
					t.Errorf("Expected status %s, got %s", test.expStatus, gotStatus)
				}

				gotColor = sanitizeColorCode(gotColor)
				if gotColor != test.expColor {
					t.Errorf("Expected color code %s, got %s", test.expColor, gotColor)
				}
			}
		})
	}
}

func sanitizeColorCode(colorCode string) string {
	newCode := strings.ReplaceAll(colorCode, "\033[", "")
	newCode = strings.ReplaceAll(newCode, "[", "")
	return newCode
}
