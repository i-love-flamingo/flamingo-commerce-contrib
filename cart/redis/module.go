package redis

import (
	"flamingo.me/dingo"
	cartInfrastructure "flamingo.me/flamingo-commerce/v3/cart/infrastructure"
	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"

	"flamingo.me/flamingo-commerce-contrib/cart/redis/infrastructure"
)

type (
	// Module for a cart storage using redis
	Module struct {
		enabled bool
	}
)

func (m *Module) Inject(
	config *struct {
		Enabled bool `inject:"config:commerce.contrib.cart.redis.enabled"`
	}) {
	if config != nil {
		m.enabled = config.Enabled
	}
}

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	if m.enabled {
		injector.Bind(new(infrastructure.CartSerializer)).To(new(infrastructure.GobSerializer))
		injector.Override(new(cartInfrastructure.CartStorage), "").To(new(infrastructure.RedisStorage)).AsEagerSingleton()
		injector.BindMap(new(healthcheck.Status), "cart.storage.redis").To(new(infrastructure.RedisStorage))
	}
}

// CueConfig defines the cart module configuration
func (*Module) CueConfig() string {
	// language=cue
	return `
commerce: {
	contrib: {
		cart: {
			redis: {
				enabled: bool | *true
				keyPrefix: string | *"cart:"
				ttl: {
					guest: string | *"24h"
					customer: string | *"168h"
				}
				address: string | *"127.0.0.1:6379"
				network: "unix" | *"tcp"
				username?: string & !=""
				password: string | *""
				idleConnections: number | *10
				database: float | int | *0
				tls: bool | *false
			}
		}
	}
}`
}
