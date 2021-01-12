package reporter

import (
	"context"
	"testing"
)

func TestReport_should_request_charger_state_from_driver(t *testing.T) {
	fakeReportConfiguration := ReporterConfiguration{
		ChargerName:       "",
		ClientCertificate: "",
		ChargerType:       "",
	}

	cancel := context.Background()
	subject := NewReporter(fakeReportConfiguration, nil, cancel, nil)
	t.Log(subject)
}
