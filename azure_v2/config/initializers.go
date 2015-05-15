package config

import (
	"gopkg.in/alecthomas/kingpin.v1"
	"log"
	"log/syslog"
	"os"
)

const (
	version    = "0.0.1"
	ApiVersion = "2014-12-01-Preview"
)

var (
	app                = kingpin.New("azure", "Azure V2 RightScale Self-Service plugin.")
	ListenFlag         = app.Flag("listen", "Hostname and port to listen on, e.g. 'localhost:8080' - hostname is optional").Default(":8080").String()
	ClientIdCred       = app.Arg("client", "The client id of the application that is registered in Azure Active Directory.").Required().String()
	ClientSecretCred   = app.Arg("secret", "The client key of the application that is registered in Azure Active Directory.").Required().String()
	ResourceCred       = app.Arg("resource", "The App ID URI of the web API (secured resource).").Required().String()
	SubscriptionIdCred = app.Arg("subscription", "The client subscription id.").Required().String()
	RefreshTokenCred   = app.Arg("refresh_token", "The token used for refreshing access token.").Required().String()
	// set base url as variable to be able to modify it in the specs
	BaseUrl   		   = "https://management.azure.com"
	Logger *log.Logger // Global syslog logger
)

func init() {
	// Parse command line
	app.Version(version)
	app.Parse(os.Args[1:])

	// Initialize global syslog logger
	if l, err := syslog.NewLogger(syslog.LOG_NOTICE|syslog.LOG_LOCAL0, 0); err != nil {
		panic("azure: failed to initialize syslog logger: " + err.Error())
	} else {
		Logger = l
	}
}
