package resources

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

const (
	networkPath = "providers/Microsoft.Network/virtualNetworks"
)

type (
	networkResponseParams struct {
		ID         string      `json:"id,omitempty"`
		Name       string      `json:"name,omitempty"`
		Location   string      `json:"location"`
		Tags       interface{} `json:"tags,omitempty"`
		Etag       string      `json:"etag,omitempty"`
		Properties interface{} `json:"properties,omitempty"`
		Href       string      `json:"href,omitempty"`
	}

	networkRequestParams struct {
		Name       string                 `json:"name,omitempty"`
		Location   string                 `json:"location"`
		Properties map[string]interface{} `json:"properties,omitempty"`
	}
	networkCreateParams struct {
		Name     string `json:"name,omitempty"`
		Location string `json:"location,omitempty"`
		Group    string `json:"group_name,omitempty"`
	}
	// Network is base struct for Azure Network resource to store input create params,
	// request create params and response params gotten from cloud.
	Network struct {
		createParams   networkCreateParams
		requestParams  networkRequestParams
		responseParams networkResponseParams
	}
)

// SetupNetworkRoutes declares routes for IPAddress resource
func SetupNetworkRoutes(e *echo.Echo) {
	e.Get("/networks", listNetworks)
	e.Get("/networks/:id", listOneNetwork)
	e.Post("/networks", createNetwork)
	e.Delete("/networks/:id", deleteNetwork)

	//nested routes
	// group := e.Group("/resource_groups/:group_name/networks")
	// group.Get("", listNetworks)
	// group.Post("", createNetwork)
	// group.Delete("/:id", deleteNetwork)
}

func listNetworks(c *echo.Context) error {
	return List(c, new(Network))
}

func listOneNetwork(c *echo.Context) error {
	params := c.Request.Form
	network := Network{
		createParams: networkCreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
		},
	}
	return Get(c, &network)
}

func createNetwork(c *echo.Context) error {
	network := new(Network)
	return Create(c, network)
}

func deleteNetwork(c *echo.Context) error {
	params := c.Request.Form
	network := Network{
		createParams: networkCreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
		},
	}
	return Delete(c, &network)
}

// GetRequestParams prepares parameters for create network request to the cloud
func (n *Network) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&n.createParams)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}

	var subnets []map[string]interface{}
	n.requestParams.Name = n.createParams.Name
	n.requestParams.Location = n.createParams.Location
	n.requestParams.Properties = map[string]interface{}{
		"addressSpace": map[string]interface{}{
			"addressPrefixes": []string{"10.0.0.0/16"},
		},
		"subnets": append(subnets, map[string]interface{}{
			"name": n.createParams.Name,
			"properties": map[string]interface{}{
				"addressPrefix": "10.0.0.0/16",
			},
		}),
	}

	return n.requestParams, nil
}

// GetResponseParams is accessor function for getting access to responseParams struct
func (n *Network) GetResponseParams() interface{} {
	return n.responseParams
}

// GetPath returns full path to the sigle network
func (n *Network) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, n.createParams.Group, networkPath, n.createParams.Name, config.APIVersion)
}

// GetCollectionPath returns full path to the collection of network
func (n *Network) GetCollectionPath(groupName string) string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, networkPath, config.APIVersion)
}

// HandleResponse manage raw cloud response
func (n *Network) HandleResponse(c *echo.Context, body []byte, actionName string) error {
	if err := json.Unmarshal(body, &n.responseParams); err != nil {
		return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	href := n.GetHref(n.createParams.Group, n.responseParams.Name)
	if actionName == "create" {
		c.Response.Header().Add("Location", href)
	} else if actionName == "get" {
		n.responseParams.Href = href
	}
	return nil
}

// GetContentType returns network content type
func (n *Network) GetContentType() string {
	return "vnd.rightscale.network+json"
}

// GetHref returns network href
func (n *Network) GetHref(groupName string, networkName string) string {
	return fmt.Sprintf("/networks/%s?group_name=%s", networkName, groupName)
}
