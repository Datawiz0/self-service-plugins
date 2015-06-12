package main

import (
	"log"

	"github.com/labstack/echo"
	em "github.com/labstack/echo/middleware"
	"github.com/rightscale/go_middleware"

	// load app files
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	am "github.com/rightscale/self-service-plugins/azure_v2/middleware"
	"github.com/rightscale/self-service-plugins/azure_v2/resources"
)

func main() {
	// Serve
	s := HttpServer()
	log.Printf("Azure plugin - listening on %s under %s environment\n", *config.ListenFlag, *config.Env)
	s.Run(*config.ListenFlag)
}

// Factory method for application
// Makes it possible to do integration testing.
func HttpServer() *echo.Echo {

	// Setup middleware
	e := echo.New()
	e.Use(middleware.RequestID)                 // Put that first so loggers can log request id
	e.Use(em.Logger())                          // Log to console
	e.Use(middleware.HttpLogger(config.Logger)) // Log to syslog
	e.Use(am.AzureClientInitializer())
	e.Use(em.Recover())

	if config.DebugMode {
		e.SetDebug(true)
	}

	e.SetHTTPErrorHandler(AzureErrorHandler(e)) // override default error handler

	// Setup routes
	prefix := e.Group("/azure_plugin") // added prefix to use multiple nginx location on one SS box
	resources.SetupSubscriptionRoutes(prefix)
	resources.SetupInstanceRoutes(prefix)
	resources.SetupGroupsRoutes(prefix)
	resources.SetupStorageAccountsRoutes(prefix)
	resources.SetupProviderRoutes(prefix)
	resources.SetupNetworkRoutes(prefix)
	resources.SetupSubnetsRoutes(prefix)
	resources.SetupIpAddressesRoutes(prefix)
	resources.SetupAuthRoutes(prefix)
	resources.SetupNetworkInterfacesRoutes(prefix)
	resources.SetupImageRoutes(prefix)
	resources.SetupOperationRoutes(prefix)

	return e
}
