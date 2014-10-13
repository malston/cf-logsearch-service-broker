package logstash

import (
	"path"
)

type Instance struct {
	Basepath string
}

func (instance Instance) CommandArgs() []string {
	configFile := instance.ConfigPath()
	return []string{
		"agent", "--verbose",
		"-f", configFile,
	}
}

func (instance Instance) PidFilePath() string {
	return path.Join(instance.baseDir(), "logstash.pid")
}

func (instance Instance) ConfigPath() string {
	return path.Join(instance.baseDir(), "sample.conf")
}

func (instance Instance) baseDir() string {
	return instance.Basepath
}
