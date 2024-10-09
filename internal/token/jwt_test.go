package token

import (
	"kg/procurement/cmd/config"
	"testing"
)

func Test_newJWTManager(t *testing.T) {
	_ = newJWTManager(config.Token{}, nil)
}
