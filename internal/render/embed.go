package render

import "embed"

// Public contains public static web server content
//
//go:embed public
var Public embed.FS
