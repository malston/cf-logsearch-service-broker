package main

import (
	"fmt"
	"path"

	"github.com/cloudfoundry-incubator/cf-lager"

	"github.com/malston/cf-logsearch-broker/api"
	"github.com/malston/cf-logsearch-broker/logsearch/logstash"
	"github.com/malston/cf-logsearch-broker/system"
)

func main() {
	logger := cf_lager.New("logsearch-broker")

	commandRunner := system.OSCommandRunner{
		Logger: logger,
	}

	logstashInstance := &logstash.Instance{
		Basepath: path.Join("."),
	}

	starter := logstash.NewProcessStarter(commandRunner)
	err := starter.Start(logstashInstance)
	if err != nil {
		fmt.Errorf("logstash failed to start: %s", err)
	}

	logstashBroker := api.New(&logstash.ServiceBroker{}, logger)
	logstashBroker.Run()
}
