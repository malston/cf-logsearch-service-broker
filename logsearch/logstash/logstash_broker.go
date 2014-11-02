package logstash

import (
	. "github.com/malston/cf-logsearch-broker/api"
	"github.com/malston/cf-logsearch-broker/system"
	"github.com/pivotal-golang/lager"
	"log"
	"path"
)

type LogstashServiceBroker struct {
	ProcessController    ProcessController
	ServiceConfiguration ServiceConfiguration
	InstanceRepository   InstanceRepository
	ServiceInstanceLimit int
	Logger               lager.Logger
}

func (broker *LogstashServiceBroker) GetCatalog() []Service {
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

func (broker *LogstashServiceBroker) Provision(instanceId string, params map[string]string) (string, error) {
	return broker.CreateInstance(instanceId)
}

func (broker *LogstashServiceBroker) Bind(instanceId string, bindingId string) (interface{}, error) {
	return broker.BindInstance(instanceId, bindingId)
}

func (broker *LogstashServiceBroker) Unbind(instanceId string, bindingId string) error {
	return nil
}

func (broker *LogstashServiceBroker) Deprovision(instanceId string) error {
	return nil
}

func NewServiceBroker(brokerLogger lager.Logger) *LogstashServiceBroker {
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

	return &LogstashServiceBroker{
		ServiceConfiguration: config.ServiceConfiguration,
		ProcessController: &OSProcessController{
			Starter: NewProcessStarter(commandRunner),
			Logger:  brokerLogger,
		},
		InstanceRepository:   repo,
		ServiceInstanceLimit: config.ServiceConfiguration.ServiceInstanceLimit,
		Logger:               brokerLogger,
	}
}

func (broker *LogstashServiceBroker) CreateInstance(instanceId string) (string, error) {
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

	err = broker.InstanceRepository.CreateInstanceDirectories(instance)
	if err != nil {
		return "", err
	}

	err = broker.InstanceRepository.CreateConfig(
		map[string]interface{}{"Host": instance.Host, "Port": instance.Port},
		path.Join(instance.TempatePath(), "logstash.conf.tmpl"),
		path.Join(instance.DataFilePath(), "logstash.conf"))
	if err != nil {
		return "", err
	}

	err = broker.ProcessController.StartAndWait(instance, 3.0)
	if err != nil {
		return "", err
	}

	return "http://locahost/dashboard/instances/" + instanceId, nil
}

func (broker *LogstashServiceBroker) buildInstance(instanceId string) (*Instance, error) {

	instance := &Instance{
		Id:           instanceId,
		Basepath:     path.Join(broker.ServiceConfiguration.InstanceDataDirectory, instanceId),
		LogDir:       path.Join(broker.ServiceConfiguration.InstanceLogDirectory, instanceId),
		TemplatePath: broker.ServiceConfiguration.DefaultConfigPath,
		Port:         5512,
		Host:         broker.ServiceConfiguration.Host,
	}

	return instance, nil
}

func (broker *LogstashServiceBroker) BindInstance(instanceId, bindingId string) (interface{}, error) {
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
