package api

import (
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/malston/cf-logsearch-broker/api/handlers"
	"github.com/martini-contrib/render"
)

type ServiceBroker interface {
	Services() []Service

	Provision(instanceID string, params map[string]string) error
	Deprovision(instanceID string) error

	Bind(instanceID, bindingID string) (interface{}, error)
	Unbind(instanceID, bindingID string) error
}

type Service struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Bindable        bool                   `json:"bindable"`
	Plans           []ServicePlan          `json:"plans"`
	Metadata        ServiceMetadata        `json:"metadata"`
	Tags            []string               `json:"tags"`
	DashboardClient ServiceDashboardClient `json:"dashboard_client"`
}

type ServiceDashboardClient struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	RedirectUri string `json:"redirect_uri"`
}

type ServicePlan struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Metadata    ServicePlanMetadata `json:"metadata"`
}

type ServicePlanMetadata struct {
	Bullets     []string `json:"bullets"`
	DisplayName string   `json:"displayName"`
}

type ServiceMetadata struct {
	DisplayName      string                  `json:"displayName"`
	LongDescription  string                  `json:"longDescription"`
	DocumentationUrl string                  `json:"documentationUrl"`
	SupportUrl       string                  `json:"supportUrl"`
	Listing          ServiceMetadataListing  `json:"listing"`
	Provider         ServiceMetadataProvider `json:"provider"`
}

type ServiceMetadataListing struct {
	Blurb    string `json:"blurb"`
	ImageUrl string `json:"imageUrl"`
}

type ServiceMetadataProvider struct {
	Name string `json:"name"`
}

type EmptyResponse struct{}

type ErrorResponse struct {
	Description string `json:"description"`
}

type CatalogResponse struct {
	Services []Service `json:"services"`
}

type ProvisioningResponse struct {
	DashboardURL string `json:"dashboard_url,omitempty"`
}

type BindingResponse struct {
	Credentials interface{} `json:"credentials"`
}

// Creates v2 service broker api for a given broker
func New(serviceBroker ServiceBroker) *martini.ClassicMartini {
	m := martini.Classic()
	m.Handlers(
		handlers.HandleAuthCheck(),
		render.Renderer(),
	)

	// Fetch catalog
	m.Get("/v2/catalog", func(r render.Render) {
		catalog := CatalogResponse{
			Services: serviceBroker.Services(),
		}
		r.JSON(200, catalog)
	})

	// Provision instance
	m.Put("/v2/service_instances/:instance_id", func(params martini.Params, r render.Render, req *http.Request) {
	})

	// Create binding
	m.Put("/v2/service_instances/:instance_id/service_bindings/:binding_id", func(params martini.Params, r render.Render) {
	})

	// Remove binding
	m.Delete("/v2/service_instances/:instance_id/service_bindings/:binding_id", func(params martini.Params, r render.Render) {
	})

	// Remove instance
	m.Delete("/v2/service_instances/:instance_id", func(params martini.Params, r render.Render) {
	})

	return m
}
