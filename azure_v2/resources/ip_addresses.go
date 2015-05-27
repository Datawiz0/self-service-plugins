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
	IpAddressPath = "providers/Microsoft.Network/publicIPAddresses"
)

type IpAddress struct {
	Id         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Location   string      `json:"location"`
	Tags       interface{} `json:"tags,omitempty"`
	Etag       string      `json:"etag,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
}

func SetupIpAddressesRoutes(e *echo.Echo) {
	e.Get("/ip_addresses", listIpAddresses)
	e.Post("/ip_addresses", createIpAddress)

	//nested routes
	group := e.Group("/resource_groups/:group_name/ip_addresses")
	group.Get("", listIpAddresses)
	// group.Post("", createInstance)
	// group.Delete("/:id", deleteInstance)
}

func listIpAddresses(c *echo.Context) error {
	return lib.ListResource(c, IpAddressPath)
}

func createIpAddress(c *echo.Context) error {
	postParams := c.Request.Form
	client, _ := lib.GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, postParams.Get("group_name"), IpAddressPath, postParams.Get("name"), config.ApiVersion)
	log.Printf("Create IpAddress request with params: %s\n", postParams)
	log.Printf("Create IpAddress path: %s\n", path)
	data := IpAddress{
		Location: postParams.Get("location"),
		Properties: map[string]interface{}{
			"publicIPAllocationMethod": "Dynamic",
			//"dnsSettings":   map[string]interface{}{
			//	"domainNameLabel": postParams.Get("domain_name")
			//}
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
		return lib.GenericException(fmt.Sprintf("IpAddress creation failed: %s", string(b)))
	}

	var dat *IpAddress
	if err := json.Unmarshal(b, &dat); err != nil {
		log.Fatal("Unmarshaling failed:", err)
	}

	return c.JSON(response.StatusCode, dat)
}
