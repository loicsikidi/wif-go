package all

import (
	"github.com/loicsikidi/wif-go/pkg/compiler/functions"

	// Link in all of the functions
	_ "github.com/loicsikidi/wif-go/pkg/compiler/functions/extract"
	_ "github.com/loicsikidi/wif-go/pkg/compiler/functions/join"
	_ "github.com/loicsikidi/wif-go/pkg/compiler/functions/split"
)

// Alias these methods, so that folks can import this to get all functions.
var (
	ProvideFrom = functions.ProvideFrom
	ProvideAll  = functions.ProvideAll
)
