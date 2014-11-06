package logsearch

import (
	"github.com/cloudfoundry-incubator/cf-lager"

	"github.com/malston/cf-logsearch-broker/api"
	"github.com/malston/cf-logsearch-broker/logsearch/logstash"
)

func main() {
	logger := cf_lager.New("logsearch-broker")

	logstashBroker := api.New(logstash.NewServiceBroker(logger), logger)
	logstashBroker.Run()
}
