package logstash_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLogstash(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logstash Suite")
}
