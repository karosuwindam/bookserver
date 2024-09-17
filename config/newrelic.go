package config

import "github.com/newrelic/go-agent/v3/newrelic"

func Init_newreclic() error {
	if NewRelic.LicenseKey == "" {
		return nil
	}
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(TraData.ServiceName),
		newrelic.ConfigLicense(NewRelic.LicenseKey),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigAppLogEnabled(true),
	)
	NewRelic.App = app
	return err
}
