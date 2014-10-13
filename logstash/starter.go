package logstash

import (
	"github.com/malston/cf-logsearch-broker/system"
	"github.com/pivotal-golang/lager"
)

type Starter interface {
	Start(instance *Instance) error
}

type ProcessStarter struct {
	commandRunner system.CommandRunner
	logger        lager.Logger
}

func NewProcessStarter(commmandRunner system.CommandRunner, logger lager.Logger) ProcessStarter {
	return ProcessStarter{
		commandRunner: commmandRunner,
		logger:        logger,
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
