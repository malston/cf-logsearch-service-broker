package logstash_test

import (
	"errors"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/malston/cf-logsearch-service-broker/logsearch/logstash"
	"github.com/pivotal-golang/lager/lagertest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeCommandRunner struct {
	Commands []string
}

func (fakeCommandRunner *FakeCommandRunner) Run(name string, args ...string) error {
	cmd := name + " " + strings.Join(args, " ")
	fakeCommandRunner.Commands = append(fakeCommandRunner.Commands, cmd)

	return nil
}

var _ = Describe("Starter", func() {
	var logger *lagertest.TestLogger
	var commandRunner *FakeCommandRunner
	var instance *logstash.Instance
	var isReadyFunc logstash.IsReady
	var starter *logstash.LogstashAgentStarter

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("agent-starter")
		commandRunner = &FakeCommandRunner{}
		instance = &logstash.Instance{
			Port: 6000,
			Host: "localhost",
		}
		isReadyFunc = func(address *net.TCPAddr) bool {
			return true
		}
	})
	JustBeforeEach(func() {
		starter = &logstash.LogstashAgentStarter{
			CommandRunner: commandRunner,
			IsReady:       isReadyFunc,
		}
	})
	Describe("Start a logstash agent", func() {
		Context("when the agent starts succesfully", func() {
			It("should not error", func() {
				err := starter.Start(instance, 1*time.Second)
				Expect(err).NotTo(HaveOccurred())
			})
			It("should execute the right command to start logstash", func() {
				starter.Start(instance, 1*time.Second)
				Ω(commandRunner.Commands).To(Equal([]string{
					"logstash agent --debug -f logstash.conf -l logstash.stdout.log -w " + (strconv.Itoa(runtime.NumCPU() / 2)),
				}))
			})
		})

		Context("when the agent fails to start", func() {
			BeforeEach(func() {
				isReadyFunc = func(address *net.TCPAddr) bool {
					return false
				}
			})

			It("returns the same error that the Wait returns", func() {
				connectionTimeoutErr := errors.New("timeout")
				err := starter.Start(instance, 1*time.Second)
				Ω(err).To(Equal(connectionTimeoutErr))
			})
		})
	})
})
