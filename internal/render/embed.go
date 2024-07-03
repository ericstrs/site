package render

import "embed"

// content contains static web server content
//
//go:embed public docs
var content embed.FS
