package logstash

import (
	. "github.com/malston/cf-logsearch-service-broker/api"
	"github.com/malston/cf-logsearch-service-broker/system"
	"github.com/pivotal-golang/lager"
	"log"
	"path"
	"time"
)

// logstashServiceBroker implements the api.ServiceBroker interface.
type logstashServiceBroker struct {
	ProcessStarter       ProcessStarter
	ServiceConfiguration ServiceConfiguration
	InstanceRepository   InstanceRepository
	ServiceInstanceLimit int
	Logger               lager.Logger
	FindFreePort         func() (int, error)
}

type ProcessStarter interface {
	Start(instance *Instance, timeout time.Duration) error
}

func NewServiceBroker(brokerLogger lager.Logger) *logstashServiceBroker {
	brokerConfigPath := ConfigPath()
	config, err := ParseConfig(brokerConfigPath)
	if err != nil {
		brokerLogger.Fatal("Loading config file", err, lager.Data{
			"broker-config-path": brokerConfigPath,
		})
	}

	if err = CheckConfig(config.ServiceConfiguration); err != nil {
		brokerLogger.Fatal("Checking config file", err)
	}

	repo := &FileSystemInstanceRepository{
		LogstashConf: config.ServiceConfiguration,
	}

	commandRunner := system.OSCommandRunner{
		brokerLogger,
	}

	return &logstashServiceBroker{
		ServiceConfiguration: config.ServiceConfiguration,
		ProcessStarter:       NewProcessStarter(commandRunner),
		InstanceRepository:   repo,
		ServiceInstanceLimit: config.ServiceConfiguration.ServiceInstanceLimit,
		Logger:               brokerLogger,
		FindFreePort:         system.FindFreePort,
	}
}

func (broker *logstashServiceBroker) GetCatalog() []Service {
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

func (broker *logstashServiceBroker) Provision(instanceId string, params map[string]string) (string, error) {
	log.Printf("CREATING INSTANCE--------------------------------------------------")

	instanceCount, err := broker.InstanceRepository.GetInstanceCount()
	if err != nil {
		return "", err
	}
	if instanceCount >= broker.ServiceInstanceLimit {
		return "", ServiceInstanceLimitReachedError
	}

	_, err = broker.InstanceRepository.FindById(instanceId)
	if err == nil {
		return "", ServiceInstanceAlreadyExistsError
	}

	instance, err := broker.buildInstance(instanceId)
	if err != nil {
		return "", err
	}

	err = broker.InstanceRepository.Save(instance)
	if err != nil {
		return "", err
	}

	err = broker.ProcessStarter.Start(instance, time.Duration(30)*time.Second)
	if err != nil {
		return "", err
	}

	return "http://locahost/dashboard/instances/" + instanceId, nil
}

func (broker *logstashServiceBroker) Bind(instanceId string, bindingId string) (interface{}, error) {
	log.Printf("BINDING INSTANCE--------------------------------------------------")
	instance, err := broker.InstanceRepository.FindById(instanceId)
	if err != nil {
		return nil, ServiceInstanceDoesNotExistsError
	}

	return struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}{
		Host: instance.Host,
		Port: instance.Port,
	}, nil
}

func (broker *logstashServiceBroker) Unbind(instanceId string, bindingId string) error {
	return nil
}

func (broker *logstashServiceBroker) Deprovision(instanceId string) error {
	return nil
}

func (broker *logstashServiceBroker) buildInstance(instanceId string) (*Instance, error) {
	port, err := broker.FindFreePort()
	if err != nil {
		return nil, err
	}

	instance := &Instance{
		Id:           instanceId,
		Basepath:     path.Join(broker.ServiceConfiguration.InstanceDataDirectory, instanceId),
		LogDir:       path.Join(broker.ServiceConfiguration.InstanceLogDirectory, instanceId),
		TemplatePath: broker.ServiceConfiguration.DefaultConfigPath,
		Port:         port,
		Host:         broker.ServiceConfiguration.Host,
	}

	return instance, nil
}
