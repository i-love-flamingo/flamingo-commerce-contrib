package cart_test

import (
	"testing"

	"flamingo.me/flamingo-commerce-contrib/cart"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	t.Parallel()

	if err := config.TryModules(nil, new(cart.Module)); err != nil {
		t.Error(err)
	}
}
