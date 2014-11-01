package logstash

import (
	"fmt"

	"github.com/pivotal-golang/lager"
)

type ProcessController interface {
	Starter
	StartAndWait(instance *Instance, timeout float64) error
}

type OSProcessController struct {
	Starter
	Logger lager.Logger
}

func (controller *OSProcessController) StartAndWait(instance *Instance, timeout float64) error {
	controller.Logger.Info("Starting logstash process")
	err := controller.Start(instance)
	if err != nil {
		return fmt.Errorf("logstash failed to start: %s", err)
	}

	return nil
}
