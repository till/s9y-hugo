package hugo

import "embed"

//go:embed *.tmpl
var Template embed.FS
