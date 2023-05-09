//go:build embed

package static

import "embed"

//go:embed kodata/*
var ebox embed.FS

func init() {
	IsEmbedded = true
}
