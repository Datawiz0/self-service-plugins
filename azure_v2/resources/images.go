package resources

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

const (
	computePath = "providers/Microsoft.Compute"
)

func SetupImageRoutes(e *echo.Echo) {
	e.Get("/locations", listLocations)
	e.Get("/publishers", listPublishers)
	e.Get("/offers", listOffers)
	e.Get("/skus", listSkus)
	e.Get("/versions", listVersions)
	e.Get("/images", listImages)
}

func listImages(c *echo.Context) error {
	params := c.Request.Form
	location := params.Get("location")
	publishers, err := GetPublishers(c, location)
	if err != nil {
		return err
	}
	var result []map[string]interface{}
	for _, publisher := range publishers {
		offers, _ := GetOffers(c, location, publisher["name"].(string))
		for _, offer := range offers {
			skus, _ := GetSkus(c, location, publisher["name"].(string), offer["name"].(string))
			for _, sku := range skus {
				versions, _ := GetVersions(c, location, publisher["name"].(string), offer["name"].(string), sku["name"].(string))
				result = append(result, versions...)
			}
		}
	}

	//TODO: add hrefs or use AzureResource interface
	return c.JSON(200, result)
}

func listLocations(c *echo.Context) error {
	locations, err := GetLocations(c)
	if err != nil {
		return err
	}
	return c.JSON(200, locations)
}

func GetLocations(c *echo.Context) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("%s/subscriptions/%s/locations?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, "2015-01-01")
	locations, err := GetResources(c, path)
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func listPublishers(c *echo.Context) error {
	params := c.Request.Form
	location := params.Get("location")
	var locations []map[string]interface{}
	var err error
	if location == "" {
		locations, err = GetLocations(c)
		if err != nil {
			return err
		}
	} else {
		locations = append(locations, map[string]interface{}{"name": location})
	}

	var results []map[string]interface{}
	for _, location := range locations {
		publishers, err := GetPublishers(c, location["name"].(string))
		if err != nil {
			return err
		}
		results = append(results, publishers...)
	}
	return c.JSON(200, results)
}

func GetPublishers(c *echo.Context, locationName string) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("%s/subscriptions/%s/%s/locations/%s/publishers?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, computePath, locationName, "2015-05-01-preview")
	publishers, err := GetResources(c, path)
	if err != nil {
		fmt.Printf("SKIP FOR %s because of error: %s\n", locationName, err)
		emptyArray := make([]map[string]interface{}, 0)
		return emptyArray, nil
		//return nil, err
	}

	return publishers, nil
}
func listOffers(c *echo.Context) error {
	params := c.Request.Form
	location := params.Get("location")
	publisher := params.Get("publisher")
	if location == "" || publisher == "" {
		return eh.GenericException("Please specify both params 'location' and 'publisher'.")
	}
	offers, err := GetOffers(c, location, publisher)
	if err != nil {
		return err
	}
	return c.JSON(200, offers)
}

func GetOffers(c *echo.Context, locationName string, publisherName string) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("%s/subscriptions/%s/%s/locations/%s/publishers/%s/artifacttypes/vmimage/offers?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, computePath, locationName, publisherName, "2015-05-01-preview")
	offers, err := GetResources(c, path)
	if err != nil {
		return nil, err
	}
	return offers, nil
}

func listSkus(c *echo.Context) error {
	params := c.Request.Form
	location := params.Get("location")
	publisher := params.Get("publisher")
	offer := params.Get("offer")
	if location == "" || publisher == "" || offer == "" {
		return eh.GenericException("Please specify the follwing params: 'location', 'publisher' and 'offer'.")
	}
	skus, err := GetSkus(c, location, publisher, offer)
	if err != nil {
		return err
	}
	return c.JSON(200, skus)
}

func GetSkus(c *echo.Context, locationName string, publisherName string, offerName string) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("%s/subscriptions/%s/%s/locations/%s/publishers/%s/artifacttypes/vmimage/offers/%s/skus?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, computePath, locationName, publisherName, offerName, "2015-05-01-preview")
	skus, err := GetResources(c, path)
	if err != nil {
		return nil, err
	}
	return skus, nil
}

func listVersions(c *echo.Context) error {
	params := c.Request.Form
	location := params.Get("location")
	publisher := params.Get("publisher")
	offer := params.Get("offer")
	sku := params.Get("sku")
	if location == "" || publisher == "" || offer == "" || sku == "" {
		return eh.GenericException("Please specify the follwing params: 'location', 'publisher', 'offer' and 'sku'.")
	}
	versions, err := GetVersions(c, location, publisher, offer, sku)
	if err != nil {
		return err
	}
	return c.JSON(200, versions)
}

func GetVersions(c *echo.Context, locationName string, publisherName string, offerName string, skuName string) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("%s/subscriptions/%s/%s/locations/%s/publishers/%s/artifacttypes/vmimage/offers/%s/skus/%s/versions?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, computePath, locationName, publisherName, offerName, skuName, "2015-05-01-preview")
	versions, err := GetResources(c, path)
	if err != nil {
		return nil, err
	}
	return versions, nil
}
