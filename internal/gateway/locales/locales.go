package locales

import "embed"

//go:embed en.yaml zh.yaml
var Locales embed.FS

const (
	NoPermission = "no.permission"
)
