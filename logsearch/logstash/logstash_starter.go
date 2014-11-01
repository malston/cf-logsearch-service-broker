package logstash

import (
	"github.com/malston/cf-logsearch-broker/system"
	"log"
)

type Starter interface {
	Start(instance *Instance) error
}

type ProcessStarter struct {
	commandRunner system.CommandRunner
}

func NewProcessStarter(commmandRunner system.CommandRunner) Starter {
	return ProcessStarter{
		commandRunner: commmandRunner,
	}
}

func (starter ProcessStarter) Start(instance *Instance) error {
	err := starter.commandRunner.Run("logstash", instance.CommandArgs()...)
	if err != nil {
		log.Printf("logstash failed to start: %s", err)
		return err
	}

	return nil
}
