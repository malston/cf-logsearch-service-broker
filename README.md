# cf-logsearch-service-broker 

A [Cloud Foundry](http://docs.cloudfoundry.org/services/api.html) service broker for [Logsearch](http://www.logsearch.io/).

## Install

```
go get github.com/malston/cf-logsearch-service-broker
cd $GOPATH/src/github.com/malston/cf-logsearch-service-broker && godep get
```

## Build & Run

```
go build -o bin/broker main.go
bin/broker
```

## Running tests

```
ginkgo -r --randomizeAllSpecs --failOnPending --skipMeasurements --trace --race  --cover
```

## Development

The `cf-logsearch-service-broker`
 uses [godep](https://github.com/tools/godep) to manage `go` dependencies.

All `go` packages required to run the broker are vendored into the `Godeps` directory.

When making changes to the code that requires additional `go` packages, you should use the workflow described in the
[Add or Update a Dependency](https://github.com/tools/godep#add-a-dependency) section of the godep documentation.

## Debugging

Turn on debug logging by setting the logLevel flag to `debug`
```
bin/broker -logLevel=debug
```