package main

import (
	"log"
	"os"

	"github.com/pivotal-golang/lager"

	"github.com/cloudfoundry-incubator/cf-lager"
	"github.com/logsearch/cf-logsearch-broker/broker"
	"github.com/logsearch/cf-logsearch-broker/system"
)

func main() {
	brokerConfigPath := broker.ConfigPath()
	httpLogger := log.New(os.Stdout, "", 0)
	brokerLogger := cf_lager.New("logsearch-broker")
}
