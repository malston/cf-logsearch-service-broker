package main

import (
	"github.com/cloudfoundry-incubator/cf-lager"

	"github.com/malston/cf-logsearch-service-broker/api"
	"github.com/malston/cf-logsearch-service-broker/logsearch/logstash"
)

func main() {
	logger := cf_lager.New("logsearch-broker")

	logstashBroker := api.New(logstash.NewServiceBroker(logger), logger)
	logstashBroker.Run()
}
