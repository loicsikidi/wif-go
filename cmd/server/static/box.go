package static

import "embed"

var Box embed.FS

var IsEmbedded bool

func init() {
	Box = ebox
}
