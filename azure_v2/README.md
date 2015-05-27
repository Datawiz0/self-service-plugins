#Azure plugin

Azure plugin is a Go application which serves REST HTTP requests and provides ability to iterate with Azure V2 API

##Build application

```
go build -o azure_plugin
```

## Usage

```
azure_plugin --help
usage: azure_plugin [<flags>] <client> <secret> <resource> <subscription> <refresh_token>

Azure V2 RightScale Self-Service plugin.

Flags:
  --help            Show help.
  --listen=":8080"  Hostname and port to listen on, e.g. 'localhost:8080' - hostname is optional
  --version         Show application version.

Args:
  <client>         The client id of the application that is registered in Azure Active Directory.
  <secret>         The client key of the application that is registered in Azure Active Directory.
  <resource>       The App ID URI of the web API (secured resource).
  <subscription>   The client subscription id.
  <refresh_token>  The token used for refreshing access token.
```

##New cloud registration
First step of cloud registration is registering RS application in the client Active Directory
in order to get ability to use application specific access token.
curl -v 'http://localhost:8080/application/register'

##Make requests
With no access token passed in the cookies
curl -v 'http://localhost:8080/instances'
This kind of call will go through azure oauth to get application specific access token.
Note: application should be registered in advance

Pass access token in the cookies
curl -v -b "AccessToken=eyJ0eXAiOiJKV1QiLCJhbGci..." 'http://localhost:8080/instances'
Note: could be used either user or app specific access token but take into account that plugin doesn't refresh token automatically

##Run tests

```
go test ./resources
```