package helper

import (
	"context"
	"fmt"

	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce-contrib/cart/redis"
	"flamingo.me/flamingo-commerce/v3/breadcrumbs"
	"flamingo.me/flamingo-commerce/v3/cart"
	"flamingo.me/flamingo-commerce/v3/category"
	"flamingo.me/flamingo-commerce/v3/checkout"
	"flamingo.me/flamingo-commerce/v3/customer"
	"flamingo.me/flamingo-commerce/v3/order"
	"flamingo.me/flamingo-commerce/v3/price"
	"flamingo.me/flamingo-commerce/v3/product"
	"flamingo.me/flamingo-commerce/v3/search"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	projectTestGraphql "flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/graphql"
	integrationCart "flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/modules/cart"
	fakeCustomer "flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/modules/customer"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/modules/payment"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/modules/placeorder"
	"flamingo.me/flamingo-commerce/v3/w3cdatalayer"
	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/core/auth"
	fakeAuth "flamingo.me/flamingo/v3/core/auth/fake"
	"flamingo.me/flamingo/v3/core/healthcheck"
	"flamingo.me/flamingo/v3/core/locale"
	"flamingo.me/flamingo/v3/core/requestlogger"
	"flamingo.me/flamingo/v3/core/robotstxt"
	"flamingo.me/flamingo/v3/core/security"
	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework"
	"flamingo.me/flamingo/v3/framework/cmd"
	"flamingo.me/flamingo/v3/framework/config"
	flamingoFramework "flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/prefixrouter"
	"flamingo.me/flamingo/v3/framework/web/filter"
	"flamingo.me/graphql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// modulesDemoProject return slice of modules that we want to have in our example app for testing
func modulesDemoProject() []dingo.Module {
	return []dingo.Module{
		new(framework.InitModule),
		new(cmd.Module),
		new(zap.Module),
		new(flamingoFramework.SessionModule),
		new(prefixrouter.Module),
		new(product.Module),
		new(locale.Module),
		new(customer.Module),
		new(fakeCustomer.Module),
		new(cart.Module),
		new(checkout.Module),
		new(search.Module),
		new(category.Module),
		new(requestlogger.Module),
		new(filter.DefaultCacheStrategyModule),
		new(auth.WebModule),
		new(fakeAuth.Module),
		new(breadcrumbs.Module),
		new(order.Module),
		new(healthcheck.Module),
		new(w3cdatalayer.Module),
		new(robotstxt.Module),
		new(security.Module),
		new(opencensus.Module),
		new(price.Module),
		new(projectTestGraphql.Module),
		new(graphql.Module),
		new(payment.Module),
		new(placeorder.Module),
		new(integrationCart.Module),
		new(redis.Module), // redis cart
	}
}

var redisC testcontainers.Container

func startUpDockerRedis(configMap config.Map) config.Map {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	var err error

	redisC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	port, err := redisC.MappedPort(ctx, "6379")
	if err != nil {
		panic(err)
	}

	host, err := redisC.Host(ctx)
	if err != nil {
		panic(err)
	}

	address := fmt.Sprintf("%s:%s", host, port.Port())

	if configMap == nil {
		configMap = config.Map{}
	}

	err = configMap.Add(config.Map{"commerce.contrib.cart.redis.address": address})
	if err != nil {
		panic(err)
	}

	return configMap
}

// BootupDemoProject boots up a complete demo project
func BootupDemoProject(configDir string) integrationtest.BootupInfo {
	return integrationtest.Bootup(modulesDemoProject(), configDir, startUpDockerRedis(nil))
}

// GenerateGraphQL generates the graphql interfaces for the demo project and saves to filesystem.
// use via makefile - each time you modify the schema
func GenerateGraphQL() {
	application, err := flamingo.NewApplication(modulesDemoProject(), flamingo.ConfigDir("config"))
	if err != nil {
		panic(err)
	}

	servicesI, err := application.ConfigArea().Injector.GetInstance(new([]graphql.Service))
	if err != nil {
		panic(err)
	}

	var ok bool

	services, ok := servicesI.([]graphql.Service)
	if !ok {
		panic("services not of correct type")
	}
	
	err = graphql.Generate(services, "graphql", "graphql/schema")
	if err != nil {
		panic(err)
	}
}
