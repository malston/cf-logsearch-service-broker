package logstash

import (
	"github.com/malston/cf-logsearch-broker/system"
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
	// execute : bin/logstash agent --verbose -f sample.conf
	err := starter.commandRunner.Run("logstash", instance.CommandArgs()...)
	if err != nil {
		return err
	}

	return nil
}
