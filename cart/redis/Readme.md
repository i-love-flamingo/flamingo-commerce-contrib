# Redis Cart Storage Module

This module offers an implementation of the `CartStorage` interface using Redis.

Note that by adding this module the default `InMemoryStorage` binding is overwritten and a running Redis is expected.
This can be disabled again via a configuration flag.

## Usage

### Configurations

For all possible configurations you can check the `module.go` (CueConfig function)
As always you can also dump the current configuration with the "config" Flamingo command.

Here is a typical configuration
```yaml
  commerce.contrib.cart.redis:
    # enables the module by binding the concrete implementations
    # if set to true a running redis is expected
    enabled: true
    # prefix for keys used to store carts in the redis database
    # will be suffixed with the ID of the cart to be stored
    keyPrefix: "cart:"
    # time to live for entries in the redis database, can differ for guests and logged-in customers
    ttl: 
      guest: "24h"
      customer: "168h"
    # address of the redis  
    address: "example.com:6379"
    # network type, either tcp or unix
    network: "tcp"
    # password of the redis
    password: "pass"
    # maximum number of socket connections
    idleConnections: 10
    # database to be selected of the redis
    database: 0
    # if TLS should be negotiated
    tls: false
```

### Serialization

By default, the carts are serialized using gob.
For other forms of serialization it is possible to reimplement the `CartSerializer` interface and overwrite the binding.
