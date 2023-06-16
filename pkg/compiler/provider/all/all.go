package all

import (
	"github.com/loicsikidi/wif-go/pkg/compiler/provider"

	// Register all providers
	_ "github.com/loicsikidi/wif-go/pkg/compiler/provider/oidc"
	_ "github.com/loicsikidi/wif-go/pkg/compiler/provider/saml"
)

var (
	AmbientProvider = provider.AmbientProvider
)
