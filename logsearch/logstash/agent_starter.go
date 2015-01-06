package logstash

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/malston/cf-logsearch-service-broker/system"
)

// LogstashAgentStarter implements the broker.go ProcessStarter interface.
type LogstashAgentStarter struct {
	CommandRunner system.CommandRunner
	IsReady       IsReady
}

type Action func(success chan<- struct{}, terminate <-chan struct{})

type IsReady func(address *net.TCPAddr) bool

func NewProcessStarter(commandRunner system.CommandRunner) ProcessStarter {
	return &LogstashAgentStarter{
		CommandRunner: commandRunner,
		IsReady:       isListening,
	}
}

func (starter *LogstashAgentStarter) Start(instance *Instance, timeout time.Duration) error {
	err := starter.CommandRunner.Run("logstash", instance.CommandArgs()...)
	if err != nil {
		return fmt.Errorf("logstash failed to start: %s", err)
	}
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", instance.Host, instance.Port))
	if err != nil {
		log.Printf("could not resolve address %s:%d because %s", instance.Host, instance.Port, err)
	}

	return starter.Wait(address, timeout)
}

func (starter *LogstashAgentStarter) Wait(address *net.TCPAddr, timeout time.Duration) error {
	return PerformActionWithin(timeout, func(success chan<- struct{}, terminate <-chan struct{}) {
		for {
			select {
			case <-terminate:
				return
			case <-time.After(10 * time.Millisecond):
				if starter.IsReady(address) {
					close(success)
					return
				}
			}
		}
	})
}

func isListening(address *net.TCPAddr) bool {
	_, err := net.DialTCP("tcp", nil, address)
	return err == nil
}

func PerformActionWithin(timeout time.Duration, action Action) error {
	success := make(chan struct{})
	terminate := make(chan struct{})
	go action(success, terminate)
	select {
	case <-success:
		return nil
	case <-time.After(timeout):
		close(terminate)
		return errors.New("timeout")
	}
}
