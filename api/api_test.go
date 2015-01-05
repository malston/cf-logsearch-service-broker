package api_test

import (
	"net/http"
	"os"

	. "github.com/malston/cf-logsearch-service-broker/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeServiceBroker struct {
	ServiceBroker
}

func (fsb *FakeServiceBroker) GetCatalog() []Service {
	return []Service{
		Service{
			Id:          "124b3b9f-89b5-4ee0-b299-850a47c4a30d",
			Name:        "logsearch-service",
			Description: "Logsearch Service for Cloud Foundry v2",
			Bindable:    true,
			DashboardClient: DashboardClient{
				Id:          "logsearch-service-client",
				Secret:      "s3cr3t",
				RedirectUri: "https://dashboard.com",
			},
			Plans: []Plan{
				Plan{
					Id:          "dc851bfa-b23c-4e07-ae4d-26a5c403ce97",
					Name:        "default",
					Description: "The default Logsearch plan",
					Metadata: PlanMetadata{
						Bullets:     []string{},
						DisplayName: "Logsearch",
					},
				},
			},
			Metadata: ServiceMetadata{
				DisplayName:      "Logsearch",
				LongDescription:  "Logsearch is open source software from City Index used to search logs using the power of ELK",
				DocumentationUrl: "http://documentation.com",
				SupportUrl:       "http://support.com",
				Listing: ServiceMetadataListing{
					Blurb:    "Logsearch ...",
					ImageUrl: "http://image.com/image.png",
				},
				Provider: ServiceMetadataProvider{
					Name: "Logsearch.io",
				},
			},
			Tags: []string{
				"logging",
				"logsearch",
			},
		},
	}
}

func (fsb *FakeServiceBroker) Provision(instanceId string, _ map[string]string) (string, error) {
	return "http://locahost/dashboard/instances/" + instanceId, nil
}

var _ = Describe("service broker api", func() {
	var (
		fakeServiceBroker *FakeServiceBroker
	)
	Describe("a logsearch catalog", func() {
		BeforeEach(func() {
			fakeServiceBroker = new(FakeServiceBroker)
			os.Setenv("LOGSEARCH_BROKER_USERNAME", "username")
			os.Setenv("LOGSEARCH_BROKER_PASSWORD", "password")
		})
		AfterEach(func() {
			os.Setenv("LOGSEARCH_BROKER_USERNAME", "")
			os.Setenv("LOGSEARCH_BROKER_PASSWORD", "")
		})
		Context("when catalog is fetched with valid credentials", func() {
			It("returns a 200 status code", func() {
				response := AuthorizedRequest("GET", "/v2/catalog", fakeServiceBroker)
				Expect(response.Code).To(Equal(200))
			})
			It("returns valid catalog json", func() {
				response := AuthorizedRequest("GET", "/v2/catalog", fakeServiceBroker)
				Expect(response.Body).To(MatchJSON(Fixture("catalog.json")))
			})
		})
		Context("when catalog is fetched with invalid credentials", func() {
			It("returns a 401 status code", func() {
				response := Request("GET", "/v2/catalog", "badusername", "badpassword", fakeServiceBroker)
				Expect(response.Code).To(Equal(http.StatusUnauthorized))
			})
		})
	})
})
