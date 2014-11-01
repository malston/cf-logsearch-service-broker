package logstash

import (
	"errors"
	"fmt"
	"os"

	"github.com/fraenkel/candiedyaml"
)

type ServiceConfiguration struct {
	Host                  string            `yaml:"host"`
	DefaultConfigPath     string            `yaml:"conf_path"`
	InstanceDataDirectory string            `yaml:"data_directory"`
	InstanceLogDirectory  string            `yaml:"log_directory"`
	ServiceInstanceLimit  int               `yaml:"service_instance_limit"`
	CommandMapping        map[string]string `yaml:"command_mapping"`
}

type Config struct {
	ServiceConfiguration ServiceConfiguration `yaml:"logstash"`
}

func ParseConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := candiedyaml.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func CheckConfig(config ServiceConfiguration) error {
	err := checkPathExists(config.DefaultConfigPath, "Logstash DefaultConfigPath")
	if err != nil {
		return err
	}

	err = checkPathExists(config.InstanceDataDirectory, "Logstash InstanceDataDirectory")
	if err != nil {
		return err
	}

	err = checkPathExists(config.InstanceLogDirectory, "Logstash InstanceLogDirectory")
	if err != nil {
		return err
	}

	return nil
}

func checkPathExists(path string, description string) error {
	_, err := os.Stat(path)
	if err != nil {
		errMessage := fmt.Sprintf(
			"File '%s' (%s) not found",
			path,
			description)
		return errors.New(errMessage)
	}
	return nil
}

func ConfigPath() string {
	brokerConfigYamlPath := os.Getenv("BROKER_CONFIG_PATH")
	if brokerConfigYamlPath == "" {
		brokerConfigYamlPath = "logsearch/logstash/assets/logstash_config.yml"
		// panic("BROKER_CONFIG_PATH not set")
	}
	return brokerConfigYamlPath
}
