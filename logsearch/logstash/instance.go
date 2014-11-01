package logstash

import (
	"path"
)

type Instance struct {
	Id       string
	Basepath string
	ConfPath string
	LogDir   string
	Host     string
	Port     int
}

func (instance Instance) CommandArgs() []string {
	// port := strconv.Itoa(instance.Port)
	return []string{
		"agent",
		"--debug",
		"-f", instance.ConfigPath(),
		"-l", instance.LogFilePath(),
		"-w", "2",
		// "-a", instance.Host,
		// "-p", port,
		// "--pidfile", instance.PidFilePath(),
		// ">>" + instance.LogFilePath(),
		// "2>>" + instance.LogFilePath(),
	}
}

func (instance Instance) ConfigPath() string {
	return path.Join(instance.ConfPath, "logstash.conf")
}

func (instance Instance) LogFilePath() string {
	return path.Join(instance.LogDir, "logstash.stdout.log")
}

func (instance Instance) baseDir() string {
	return instance.Basepath
}
