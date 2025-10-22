package static

import (
	"embed"
)

//go:embed all:public
var PublicFS embed.FS
