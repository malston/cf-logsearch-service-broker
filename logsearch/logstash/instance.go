package logstash

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
)

type Instance struct {
	Id           string
	Basepath     string
	LogDir       string
	Host         string
	Port         int
	TemplatePath string
}

func (instance Instance) CommandArgs() []string {
	return []string{
		"agent",
		"--debug",
		"-f", instance.ConfigPath(),
		"-l", instance.LogFilePath(),
		"-w", strconv.Itoa(runtime.NumCPU() / 2),
	}
}

func (instance Instance) Address() string {
	return fmt.Sprintf("%s:%d", instance.Host, instance.Port)
}

func (instance Instance) ConfigPath() string {
	return path.Join(instance.baseDir(), "logstash.conf")
}

func (instance Instance) LogFilePath() string {
	return path.Join(instance.LogDir, "logstash.stdout.log")
}

func (instance Instance) DataFilePath() string {
	return instance.baseDir()
}

func (instance Instance) TempatePath() string {
	return instance.TemplatePath
}

func (instance Instance) baseDir() string {
	return instance.Basepath
}
