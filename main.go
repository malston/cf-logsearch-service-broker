package main

import (
	"fmt"
	"path"

	"github.com/cloudfoundry-incubator/cf-lager"
	"github.com/pivotal-golang/lager"

	"github.com/malston/cf-logsearch-broker/logstash"
	"github.com/malston/cf-logsearch-broker/system"
)

func main() {
	var logger lager.Logger
	logger = cf_lager.New("logsearch-broker")

	commandRunner := system.OSCommandRunner{
		Logger: logger,
	}

	instance := &logstash.Instance{
		Basepath: path.Join("."),
	}

	starter := logstash.NewProcessStarter(commandRunner, logger)
	err := starter.Start(instance)
	if err != nil {
		fmt.Errorf("logstash failed to start: %s", err)
	}
}
