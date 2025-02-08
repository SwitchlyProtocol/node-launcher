package monitor

import (
	"public-alerts/internal/notify"
	"testing"
)

// Mock implementation of TorManipulationMonitor for testing
type MockTorManipulationMonitor struct {
	conditionMet bool
}

func (m *MockTorManipulationMonitor) Check() ([]notify.Alert, error) {
	alert := notify.Alert{
		Message: `**🚨 ALERT: TOR PRICE 🚨**
	> TOR Price in USD: $4.000
	> TOR Price Delta: 300.000%
	> ---------------------------
	> Rune Price in TOR: 0.500
	> Rune Price in USD: $2.000
	> Rune Price Delta: -75.000%
	> Check TOR Anchor pool depths and median value.`,
	}

	if m.conditionMet {
		return []notify.Alert{alert}, nil
	}

	return nil, nil
}

func TestCheckTorManipulation_NoAlert(t *testing.T) {
	monitor := &MockTorManipulationMonitor{conditionMet: false}

	alerts, err := monitor.Check()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(alerts) != 0 {
		t.Fatalf("expected no alerts, got some")
	}
}

func TestCheckTorManipulation_WithAlert(t *testing.T) {
	monitor := &MockTorManipulationMonitor{conditionMet: true}

	alerts, err := monitor.Check()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(alerts) != 1 {
		t.Fatalf("expected one alert, got %d", len(alerts))
	}

	expectedMsg := `**🚨 ALERT: TOR PRICE 🚨**
	> TOR Price in USD: $4.000
	> TOR Price Delta: 300.000%
	> ---------------------------
	> Rune Price in TOR: 0.500
	> Rune Price in USD: $2.000
	> Rune Price Delta: -75.000%
	> Check TOR Anchor pool depths and median value.`

	if alerts[0].Message != expectedMsg {
		t.Errorf("unexpected alert message: %s", alerts[0].Message)
	}
}
