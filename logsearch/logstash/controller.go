package logstash

import (
	"fmt"

	"github.com/pivotal-golang/lager"
)

// logstashProcessController implements the ProcessController interface.
type logstashProcessController struct {
	ProcessStarter
	Logger lager.Logger
}

func (controller *logstashProcessController) StartAndWait(instance *Instance, timeout float64) error {
	controller.Logger.Info("Starting logstash process")
	err := controller.Start(instance)
	if err != nil {
		return fmt.Errorf("logstash failed to start: %s", err)
	}

	return nil
}
