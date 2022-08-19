package interactions

import (
	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/rest"
)

type (
	// Config contains the configuration values for the interactions App
	Config struct {
		// PublicKey is your interactions public key provided on the Discord developers site
		PublicKey string
		// Logger allows you to specify a custom logrus Logger for the App to use
		Logger *zerolog.Logger
		// Token (optional) is your Discord token that will be passed to the internal REST client
		Token string
		// RESTClient (optional) is the REST client you are overriding with. Useful for proxies.
		RESTClient *rest.Client
	}
)
