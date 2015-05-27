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
	virtualMachinesPath = "providers/Microsoft.Compute/virtualMachines"
)

type Instance struct {
	ProvisioningState interface{}            `json:"provisioningState,omitempty"`
	InstanceView      interface{}            `json:"instanceView,omitempty"`
	HardwareProfile   interface{}            `json:"hardwareProfile,omitempty"`
	NetworkProfile    interface{}            `json:"networkProfile,omitempty"`
	StorageProfile    interface{}            `json:"storageProfile,omitempty"`
	Id                string                 `json:"id,omitempty"`
	Name              string                 `json:"name"`
	Type              string                 `json:"type,omitempty"`
	Location          string                 `json:"location"`
	Properties        map[string]interface{} `json:"properties,omitempty"` // used for create instance
}

func SetupInstanceRoutes(e *echo.Echo) {
	//get all instances from all groups
	e.Get("/instances", listInstances)

	//nested routes
	group := e.Group("/resource_groups/:group_name/instances")
	group.Get("", listInstances)
	group.Post("", createInstance)
	group.Delete("/:id", deleteInstance)
}

func listInstances(c *echo.Context) error {
	group_name := c.Param("group_name")
	if group_name != "" {
		path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, virtualMachinesPath, config.ApiVersion)
		instances, err := lib.GetResources(c, path)
		if err != nil {
			return err
		}
		return c.JSON(200, instances)
	} else {
		path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, "2015-01-01")
		relativePath := fmt.Sprintf("/%s?api-version=%s", virtualMachinesPath, config.ApiVersion)
		instances, err := lib.ListNestedResources(c, path, relativePath)
		if err != nil {
			return err
		}
		return c.JSON(200, instances)
	}

}

func deleteInstance(c *echo.Context) error {
	group_name := c.Param("group_name")
	if group_name == "" {
		return lib.GenericException("Parameter 'group_name' is required.")
	}
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, virtualMachinesPath, c.Param("id"), config.ApiVersion)
	return lib.DeleteResource(c, path)
}

func getInstances(c *echo.Context, group_name string) (int, []interface{}) {
	client, _ := lib.GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, virtualMachinesPath, config.ApiVersion)
	log.Printf("Get Instances request: %s\n", path)
	resp, err := client.Get(path)

	if err != nil {
		log.Fatal("Get:", err)
	}
	defer resp.Body.Close()
	var m map[string][]interface{}
	b, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(b, &m)

	return resp.StatusCode, m["value"]
}

// check out that provider is already registered - https://msdn.microsoft.com/en-us/library/azure/dn790548.aspx
func createInstance(c *echo.Context) error {
	postParams := c.Request.Form
	client, _ := lib.GetAzureClient(c)
	var networkInterfaces []map[string]interface{}
	instanceParams := Instance{
		Name:     postParams.Get("name"),
		Location: postParams.Get("location"),
		Properties: map[string]interface{}{
			"hardwareProfile": map[string]interface{}{"vmSize": postParams.Get("instance_type_uid")},
			"storageProfile": map[string]interface{}{
				"osDisk": map[string]interface{}{
					"vhd": map[string]interface{}{
						"uri": "https://khrvitestgo.blob.core.windows.net/vhds/khrvi_image-os-2015-05-18.vhd"},
					"name":   "os-" + postParams.Get("name") + "-rs",
					"osType": "Linux"},
				// "sourceImage": map[string]interface{}{
				// 	"id": "/2d2b2267-ff0a-46d3-9912-8577acb18a0a/services/images/7bb63e06fb004b2597e854325d2fe7b9__Test-Windows-Server-2012-Datacenter-201401.01-en.us-127GB.vhd",
				// },
				// "destinationVhdsContainer": "http://khrvitestgo.blob.core.windows.net/vhds", // hard coded for now...should be used Placement group
			},
			"networkProfile": map[string]interface{}{
				"networkInterfaces": append(networkInterfaces, map[string]interface{}{
					"id": "/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/khrvi_ni",
				}),
			},
		},
	}

	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, c.Param("group_name"), virtualMachinesPath, instanceParams.Name, config.ApiVersion)
	log.Printf("Create Instances request with params: %s\n", postParams)
	log.Printf("Create Instances path: %s\n", path)

	by, err := json.Marshal(instanceParams)
	var reader io.Reader
	reader = bytes.NewBufferString(string(by))
	log.Printf("READER: %s", reader)
	request, _ := http.NewRequest("PUT", path, reader)
	request.Header.Add("Content-Type", config.MediaType)
	request.Header.Add("Accept", config.MediaType)
	request.Header.Add("User-Agent", config.UserAgent)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("Post:", err)
	}
	defer response.Body.Close()
	b, _ := ioutil.ReadAll(response.Body)
	return c.JSON(response.StatusCode, string(b))
}
