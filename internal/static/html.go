package static

import (
	"embed"
)

var (
	//go:embed templates/**.html
	HandlerTemplates embed.FS
)
