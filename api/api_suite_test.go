package api_test

import (
	"fmt"
	"github.com/malston/cf-logsearch-broker/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
)

var (
	response *httptest.ResponseRecorder
)

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CF Service Broker v2 API Suite")
}

func Request(method string, route string, broker api.ServiceBroker) {
	m := api.New(broker)
	request, _ := http.NewRequest(method, route, nil)
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
}

func Fixture(name string) string {
	filePath := path.Join("fixtures", name)
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("Could not read fixture: %s", name))
	}

	return string(contents)
}
