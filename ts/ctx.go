package ts

import (
	"html/template"
	"time"
)

func (ts *TapeStatsApp) templateContextNow() string {
	return time.Now().UTC().String()
}

func (ts *TapeStatsApp) templateContextNowYear() int {
	return time.Now().UTC().Year()
}

func (ts *TapeStatsApp) GetTemplateContext() template.FuncMap {
	return template.FuncMap{
		"now":     ts.templateContextNow,
		"nowYear": ts.templateContextNowYear,
	}
}
