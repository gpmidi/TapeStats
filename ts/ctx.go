package ts

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"time"
)

func (ts *TapeStatsApp) templateContextNow() string {
	return time.Now().UTC().String()
}

func (ts *TapeStatsApp) templateContextNowYear() string {
	return fmt.Sprint(time.Now().UTC().Year())
}

func (ts *TapeStatsApp) SetTemplateContext(r *gin.Engine) {
	r.SetFuncMap(template.FuncMap{
		"now":     ts.templateContextNow,
		"nowYear": ts.templateContextNowYear,
	})
}
