package middleware

import (
	"testing"
)

func Test_newAuthMiddleware(t *testing.T) {
	_ = NewAuthMiddleware(nil)
}
