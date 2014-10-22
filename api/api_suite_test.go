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

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CF Service Broker v2 API Suite")
}

func Request(method string, route string, username string, password string, broker api.ServiceBroker) *httptest.ResponseRecorder {
	return makeRequest(method, route, username, password, broker)
}

func AuthorizedRequest(method string, route string, broker api.ServiceBroker) *httptest.ResponseRecorder {
	return makeRequest(method, route, "username", "password", broker)
}

func UnauthorizedRequest(method string, route string, broker api.ServiceBroker) *httptest.ResponseRecorder {
	return makeRequest(method, route, "", "", broker)
}

func Fixture(name string) string {
	filePath := path.Join("fixtures", name)
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("Could not read fixture: %s", name))
	}

	return string(contents)
}

func makeRequest(method string, route string, username string, password string, broker api.ServiceBroker) *httptest.ResponseRecorder {
	m := api.New(broker)
	request, _ := http.NewRequest(method, route, nil)
	if username != "" {
		request.SetBasicAuth(username, password)
	}
	response := httptest.NewRecorder()
	m.ServeHTTP(response, request)
	return response
}
