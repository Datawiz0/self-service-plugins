package resources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

const (
	subscriptionsPath = "subscriptions"
)

// Subscription is base struct for Azure Subscription resource
type Subscription struct {
	ID             string      `json:"id"`
	Name           string      `json:"displayName"`
	State          string      `json:"state"`
	SubscriptionID string      `json:"subscriptionId"`
	Policies       interface{} `json:"subscriptionPolicies"`
}

// SetupSubscriptionRoutes declares routes for Subscription resource
func SetupSubscriptionRoutes(e *echo.Echo) {
	// e.Get("/subscriptions", listSubscriptions)
	// get a current subscription
	e.Get("/subscription", getSubscription)
}

func listSubscriptions(c *echo.Context) error {
	client, _ := GetAzureClient(c)
	path := fmt.Sprintf("%s/%s?api-version=%s", config.BaseURL, subscriptionsPath, config.APIVersion)
	log.Printf("Get Subscriptions request: %s\n", path)
	resp, err := client.Get(path)
	if err != nil {
		return eh.GenericException(fmt.Sprintf("Error has occurred while getting subscriptions: %v", err))
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var dat map[string][]*Subscription
	if err := json.Unmarshal(body, &dat); err != nil {
		return eh.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}

	return c.JSON(resp.StatusCode, dat["value"])
}

// getSubscription return info about subscription provided in creds
func getSubscription(c *echo.Context) error {
	client, _ := GetAzureClient(c)
	path := fmt.Sprintf("%s/%s/%s?api-version=%s", config.BaseURL, subscriptionsPath, *config.SubscriptionIDCred, "2015-01-01")
	log.Printf("Get Subscription request: %s\n", path)
	resp, err := client.Get(path)
	if err != nil {
		return eh.GenericException(fmt.Sprintf("Error has occurred while getting subscription: %v", err))
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var dat *Subscription
	if err := json.Unmarshal(body, &dat); err != nil {
		return eh.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}

	return c.JSON(resp.StatusCode, dat)
}
