package ts

import (
	"github.com/spf13/viper"
	"html/template"
	"time"
)

func (ts *TapeStatsApp) templateContextNow() string {
	return time.Now().UTC().String()
}

func (ts *TapeStatsApp) templateContextNowYear() int {
	return time.Now().UTC().Year()
}

func (ts *TapeStatsApp) templateContextGoogleMeasureId() string {
	return viper.GetString("google.measure.id")
}

func (ts *TapeStatsApp) GetTemplateContext() template.FuncMap {
	return template.FuncMap{
		"now":             ts.templateContextNow,
		"nowYear":         ts.templateContextNowYear,
		"googleMeasureId": ts.templateContextGoogleMeasureId,
	}
}
