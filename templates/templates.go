package views

import (
	"embed"
)

//go:embed layout/*.html *.html
var TmplFS embed.FS
