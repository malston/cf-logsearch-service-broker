package api

import (
	"errors"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/malston/cf-logsearch-broker/api/handlers"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/pivotal-golang/lager"
)

// Implements the Cloud Foundry Service Broker API
// http://docs.cloudfoundry.org/services/api.html#api-overview
type ServiceBroker interface {
	// Fetches the service catalog for the developer to select from the marketplace
	// http://docs.cloudfoundry.org/services/api.html#catalog-mgmt
	GetCatalog() []Service

	// Creates a new service resource for the developer
	// http://docs.cloudfoundry.org/services/api.html#provisioning
	Provision(instanceId string, params map[string]string) (string, error)

	// Creates a binding to a provisioned service instance for an application to use for connecting to the instance
	// http://docs.cloudfoundry.org/services/api.html#binding
	Bind(instanceId, bindingId string) (interface{}, error)

	// Removes a service instance binding so applications can no longer bind to that instance
	// http://docs.cloudfoundry.org/services/api.html#unbinding
	Unbind(instanceId, bindingId string) error

	// Deletes a provisioned service instance completely from a space so users can no longer use it
	// http://docs.cloudfoundry.org/services/api.html#deprovisioning
	Deprovision(instanceId string) error
}

// Broker API Errors
var (
	// 409 HTTP status code should be returned if the requested service instance already exists.
	ServiceInstanceAlreadyExistsError = errors.New("service instance already exists")
	// 500 HTTP status code should be returned if the instance limit for this service has been reached.
	ServiceInstanceLimitReachedError = errors.New("instance limit for this service has been reached")
)

type Service struct {
	Id              string          `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Bindable        bool            `json:"bindable"`
	Plans           []Plan          `json:"plans"`
	Metadata        ServiceMetadata `json:"metadata"`
	Tags            []string        `json:"tags"`
	DashboardClient DashboardClient `json:"dashboard_client"`
}

type DashboardClient struct {
	Id          string `json:"id"`
	Secret      string `json:"secret"`
	RedirectUri string `json:"redirect_uri"`
}

type Plan struct {
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Metadata    PlanMetadata `json:"metadata"`
}

type PlanMetadata struct {
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

type ProvisionRequest struct {
	ServiceId        string `json:"service_id"`
	PlanId           string `json:"plan_id"`
	OrganizationGuid string `json:"organization_guid"`
	SpaceGuid        string `json:"space_guid"`
}

type EmptyResponse struct{}

type ErrorResponse struct {
	Description string `json:"description"`
}

type CatalogResponse struct {
	Services []Service `json:"services"`
}

type ProvisionResponse struct {
	DashboardUrl string `json:"dashboard_url,omitempty"`
}

type BindingResponse struct {
	Credentials interface{} `json:"credentials"`
}

// Creates v2 service broker api for a given broker
func New(serviceBroker ServiceBroker, logger lager.Logger) *martini.ClassicMartini {
	m := martini.Classic()
	m.Handlers(
		handlers.HandleAuthCheck(),
		render.Renderer(),
	)

	// Fetch catalog
	m.Get("/v2/catalog", func(r render.Render) {
		catalog := CatalogResponse{
			Services: serviceBroker.GetCatalog(),
		}
		r.JSON(200, catalog)
	})

	// Provision instance
	m.Put("/v2/service_instances/:instance_id", binding.Bind(ProvisionRequest{}), func(provisionRequest ProvisionRequest, params martini.Params, r render.Render, req *http.Request) {
		instanceId := params["instance_id"]
		url, err := serviceBroker.Provision(instanceId, map[string]string{
			"organization_guid": provisionRequest.OrganizationGuid,
			"plan_id":           provisionRequest.PlanId,
			"service_id":        provisionRequest.ServiceId,
			"space_guid":        provisionRequest.SpaceGuid,
		})

		if err != nil {
			status, response := handleServiceError(err, logger)
			r.JSON(status, response)
		}

		r.JSON(201, ProvisionResponse{
			DashboardUrl: url,
		})
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

func handleServiceError(err error, logger lager.Logger) (int, interface{}) {
	logger.Error("service-broker-error", err, lager.Data{"error": err.Error()})

	switch err {
	case ServiceInstanceAlreadyExistsError:
		logger.Error("service-instance-already-exists", err)
		return 409, EmptyResponse{}
	case ServiceInstanceLimitReachedError:
		logger.Error("service-instance-limit-reached", err)
		return 500, ErrorResponse{
			Description: err.Error(),
		}
	default:
		logger.Error("unknown-error", err)
		return 500, ErrorResponse{
			Description: "an unexpected error occurred",
		}
	}
}
