package api_test

import (
	. "github.com/malston/cf-logsearch-broker/api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeServiceBroker struct {
	ServiceBroker
}

func (fsb *FakeServiceBroker) Services() []Service {
	return []Service{
		Service{
			ID:          "124b3b9f-89b5-4ee0-b299-850a47c4a30d",
			Name:        "p-logsearch-dev",
			Description: "Logsearch Service for Cloud Foundry v2",
			Bindable:    true,
			Plans: []ServicePlan{
				ServicePlan{
					ID:          "dc851bfa-b23c-4e07-ae4d-26a5c403ce97",
					Name:        "default",
					Description: "The default Logsearch plan",
					Metadata: ServicePlanMetadata{
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

var _ = Describe("service broker api", func() {
	var (
		fakeServiceBroker *FakeServiceBroker
	)
	Describe("fetching catalog", func() {
		Context("when service is available", func() {
			BeforeEach(func() {
				fakeServiceBroker = &FakeServiceBroker{}
			})
			It("returns a 200 status code", func() {
				Request("GET", "/v2/catalog", fakeServiceBroker)
				Expect(response.Code).To(Equal(200))
			})
			It("returns valid catalog json", func() {
				Request("GET", "/v2/catalog", fakeServiceBroker)
				Expect(response.Body).To(MatchJSON(Fixture("catalog.json")))
			})
		})
	})
})
