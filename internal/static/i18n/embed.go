package i18n

import (
	"embed"
)

//go:embed *.ini
var Files embed.FS

// TODO: change ini to yaml
