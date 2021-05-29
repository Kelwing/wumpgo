package interactions

import "github.com/sirupsen/logrus"

type (
	Config struct {
		PublicKey string
		Logger    *logrus.Logger
		Token     string
	}
)
