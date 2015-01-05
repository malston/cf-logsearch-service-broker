package logstash

import (
	"github.com/malston/cf-logsearch-service-broker/system"
	"log"
)

// logstashAgentStarter implements the ProcessStarter interface.
type logstashAgentStarter struct {
	commandRunner system.CommandRunner
}

func NewProcessStarter(commmandRunner system.CommandRunner) ProcessStarter {
	return logstashAgentStarter{
		commandRunner: commmandRunner,
	}
}

func (starter logstashAgentStarter) Start(instance *Instance) error {
	err := starter.commandRunner.Run("logstash", instance.CommandArgs()...)
	if err != nil {
		log.Printf("logstash failed to start: %s", err)
		return err
	}

	return nil
}
