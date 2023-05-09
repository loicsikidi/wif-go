//go:build !embed

package static

import "embed"

var ebox embed.FS

func init() {
	IsEmbedded = false
}
