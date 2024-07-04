package render

import "embed"

// content contains static web server content
//
//go:embed public
var content embed.FS
