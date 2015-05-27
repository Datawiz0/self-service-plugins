package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/lib"
)

const (
	NetworkInterfacePath = "providers/Microsoft.Network/networkInterfaces"
)

type NetworkInterface struct {
	Id         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Location   string      `json:"location"`
	Tags       interface{} `json:"tags,omitempty"`
	Etag       string      `json:"etag,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
}

func SetupNetworkInterfacesRoutes(e *echo.Echo) {
	e.Get("/network_interfaces", listNetworkInterfaces)
	e.Post("/network_interfaces", createNetworkInterface)

	//nested routes
	group := e.Group("/resource_groups/:group_name/network_interfaces")
	group.Get("", listNetworkInterfaces)
	// group.Post("", createInstance)
	// group.Delete("/:id", deleteInstance)
}

func listNetworkInterfaces(c *echo.Context) error {
	return lib.ListResource(c, NetworkInterfacePath)
}

func createNetworkInterface(c *echo.Context) error {
	postParams := c.Request.Form
	client, _ := lib.GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, postParams.Get("group_name"), NetworkInterfacePath, postParams.Get("name"), config.ApiVersion)
	log.Printf("Create NetworkInterface request with params: %s\n", postParams)
	log.Printf("Create NetworkInterface path: %s\n", path)
	var configs []map[string]interface{}
	data := NetworkInterface{
		Location: postParams.Get("location"),
		Properties: map[string]interface{}{
			"ipConfigurations": append(configs, map[string]interface{}{
				"name": postParams.Get("name") + "_ip",
				"properties": map[string]interface{}{
					"subnet": map[string]interface{}{
						"id": postParams.Get("subnet_id"),
					},
					//"privateIPAddress": "10.0.0.8",
					"privateIPAllocationMethod": "Dynamic",
					// "publicIPAddress": map[string]interface{}{
					// 	"id": ""
					// }
				},
			}),
			// "dnsSettings": map[string]interface{}{
			// 	"dnsServers": postParams.Get("dns_servers")
			// }
		},
	}

	by, err := json.Marshal(data)
	var reader io.Reader
	reader = bytes.NewBufferString(string(by))
	request, _ := http.NewRequest("PUT", path, reader)
	request.Header.Add("Content-Type", config.MediaType)
	request.Header.Add("Accept", config.MediaType)
	request.Header.Add("User-Agent", config.UserAgent)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("PUT:", err)
	}

	defer response.Body.Close()
	b, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode >= 400 {
		return lib.GenericException(fmt.Sprintf("NetworkInterface creation failed: %s", string(b)))
	}

	var dat *IpAddress
	if err := json.Unmarshal(b, &dat); err != nil {
		log.Fatal("Unmarshaling failed:", err)
	}

	return c.JSON(response.StatusCode, dat)
}
