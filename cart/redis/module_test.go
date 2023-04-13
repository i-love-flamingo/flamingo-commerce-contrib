package redis_test

import (
	"testing"

	"flamingo.me/flamingo-commerce-contrib/cart/redis"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	t.Parallel()

	if err := config.TryModules(nil, new(redis.Module)); err != nil {
		t.Error(err)
	}
}
