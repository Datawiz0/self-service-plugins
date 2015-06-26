package resources

import (
	"encoding/json"
	//"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
)

const (
	listInstancesEmptyResponse = `{"value":[]}`
	listInstancesResponse      = `{"value":[{"href":"/instances/khrvi?group_name=Group-1","id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Compute/virtualMachines/khrvi","location":"westus","name":"khrvi","properties":{"hardwareProfile":{"vmSize":"Standard_G1"},"networkProfile":{"networkInterfaces":[{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/khrvi_ni"}]},"provisioningState":"failed","storageProfile":{"dataDisks":[],"osDisk":{"caching":"ReadWrite","name":"os-asdasdasda-rs","osType":"Linux","vhd":{"uri":"https://khrvitestgo.blob.core.windows.net/vhds/khrvi_image-os-2015-05-18.vhd"}}}},"type":"Microsoft.Compute/virtualMachines"}]}`
	listOneInstanceResponse    = `{"href":"/instances/khrvi?group_name=Group-1","id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Compute/virtualMachines/khrvi","location":"westus","name":"khrvi","properties":{"hardwareProfile":{"vmSize":"Standard_G1"},"networkProfile":{"networkInterfaces":[{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/khrvi_ni"}]},"provisioningState":"failed","storageProfile":{"dataDisks":[],"osDisk":{"caching":"ReadWrite","name":"os-asdasdasda-rs","osType":"Linux","vhd":{"uri":"https://khrvitestgo.blob.core.windows.net/vhds/khrvi_image-os-2015-05-18.vhd"}}}},"type":"Microsoft.Compute/virtualMachines"}`
	recordNotFound             = `{"error":{"code":"ResourceNotFound","message":"Resource not found."}}`
)

var _ = Describe("instances", func() {

	var do *ghttp.Server
	var client *AzureClient
	var response *Response
	var err error

	BeforeEach(func() {
		do = ghttp.NewServer()
		config.BaseURL = do.URL()
		client = NewAzureClient()
	})

	AfterEach(func() {
		do.Close()
	})

	Describe("listing", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+virtualMachinesPath),
					ghttp.RespondWith(http.StatusOK, listInstancesResponse),
				),
			)
			response, err = client.Get("/resource_groups/Group-1/instances")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.instance+json;type=collection"))
		})

		It("lists all instances inside one resource group", func() {
			instances := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listInstancesResponse), &instances)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(instances["value"])
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing via 'flat' route", func() {
		BeforeEach(func() {
			subscriptionID := "test"
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups"),
					ghttp.RespondWith(http.StatusOK, `{"value": [{"name":"Group-1"}]}`),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+virtualMachinesPath),
					ghttp.RespondWith(http.StatusOK, listInstancesResponse),
				),
			)
			response, err = client.Get("/instances")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(2))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.instance+json;type=collection"))
		})

		It("lists all instances inside one resource group", func() {
			instances := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listInstancesResponse), &instances)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(instances["value"])
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing empty", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+virtualMachinesPath),
					ghttp.RespondWith(http.StatusOK, listInstancesEmptyResponse),
				),
			)
		})

		It("returns empty array", func() {
			response, err = client.Get("/resource_groups/Group-1/instances")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
			Ω(response.Body).Should(Equal("[]\n"))
		})
	})

	Describe("retrieving via 'flat' route", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+virtualMachinesPath+"/khrvi"),
					ghttp.RespondWith(http.StatusOK, listOneInstanceResponse),
				),
			)
			response, err = client.Get("/instances/khrvi?group_name=Group-1")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.instance+json"))
		})

		It("retrieves an existing instance", func() {
			var instance map[string]interface{}
			err := json.Unmarshal([]byte(listOneInstanceResponse), &instance)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(instance)
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("retrieving via 'flat' route a non-existant resource", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+virtualMachinesPath+"/khrvi1"),
					ghttp.RespondWith(http.StatusNotFound, recordNotFound),
				),
			)
		})

		It("returns 404", func() {
			response, err = client.Get("/instances/khrvi1?group_name=Group-1")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(404))
			Ω(response.Body).Should(Equal("{\"Code\":404,\"Message\":\"Could not find resource with id: khrvi1\"}\n"))
		})
	})
})
