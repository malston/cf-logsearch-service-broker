# cf-logsearch-broker 

A [Cloud Foundry](http://docs.cloudfoundry.org/services/api.html) service broker for [Logsearch](http://www.logsearch.io/).

## Install

```
go get github.com/malston/cf-logsearch-service-broker
godep get
```

## Running tests

```
ginkgo -r --randomizeAllSpecs --failOnPending --trace --race
```
