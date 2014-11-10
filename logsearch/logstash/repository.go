package logstash

import (
	"github.com/karlseguin/gerb"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

type InstanceRepository interface {
	Save(instance *Instance) error
	FindById(instanceID string) (*Instance, error)
	GetInstanceCount() (int, error)
}

type FileSystemInstanceRepository struct {
	LogstashConf ServiceConfiguration
}

type ConfigFile struct {
	Host string
	Port int
}

func (instanceRepository *FileSystemInstanceRepository) FindById(instanceId string) (*Instance, error) {
	instanceDataDir := path.Join(instanceRepository.instanceDataDirectory(), instanceId)

	_, err := os.Stat(instanceDataDir)
	if err != nil {
		return nil, err
	}

	portBytes, err := ioutil.ReadFile(path.Join(instanceDataDir, "logstash.port"))
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(strings.TrimSpace(string(portBytes)))
	if err != nil {
		return nil, err
	}

	instance := &Instance{
		Id:   instanceId,
		Port: port,
		Host: instanceRepository.LogstashConf.Host,
	}

	return instance, nil
}

func (instanceRepository *FileSystemInstanceRepository) GetInstanceCount() (int, error) {
	instances, err := instanceRepository.findAllInstances()
	return len(instances), err
}

func (instanceRepository *FileSystemInstanceRepository) Save(instance *Instance) error {
	err := instanceRepository.createBaseDirectory(instance)
	if err != nil {
		return err
	}

	err = instanceRepository.createLogDirectory(instance)
	if err != nil {
		return err
	}

	err = instanceRepository.createBindData(instance)
	if err != nil {
		return err
	}

	err = instanceRepository.createConfig(
		map[string]interface{}{"Host": instance.Host, "Port": instance.Port},
		path.Join(instance.TempatePath(), "logstash.conf.tmpl"),
		path.Join(instance.DataFilePath(), "logstash.conf"))
	if err != nil {
		return err
	}

	return nil
}

func (instanceRepository *FileSystemInstanceRepository) createBaseDirectory(instance *Instance) error {
	mkdirErr := os.MkdirAll(instance.baseDir(), 0755)
	if mkdirErr != nil {
		return mkdirErr
	}

	return nil
}

func (instanceRepository *FileSystemInstanceRepository) createLogDirectory(instance *Instance) error {
	mkdirErr := os.MkdirAll(instance.LogDir, 0755)
	if mkdirErr != nil {
		return mkdirErr
	}

	return nil
}

func (instanceRepository *FileSystemInstanceRepository) createBindData(instance *Instance) error {

	port := strconv.FormatInt(int64(instance.Port), 10)
	ioutil.WriteFile(path.Join(instance.baseDir(), "logstash.port"), []byte(port), 0644)

	return nil
}

func (instanceRepository *FileSystemInstanceRepository) instanceDataDirectory() string {
	return instanceRepository.LogstashConf.InstanceDataDirectory
}

func (instanceRepository *FileSystemInstanceRepository) instanceLogDirectory() string {
	return instanceRepository.LogstashConf.InstanceLogDirectory
}

func (instanceRepository *FileSystemInstanceRepository) findAllInstances() ([]*Instance, error) {
	instances := []*Instance{}
	log.Printf("ALL INSTANCES--------------------------------------------------")

	instanceDirs, err := ioutil.ReadDir(instanceRepository.instanceDataDirectory())
	if err != nil {
		log.Printf("ALL INSTANCES-----err reading dir %v", err)
		return instances, err
	}

	log.Printf("LOOPING ALL INSTANCES--------------------------------------------------")
	for _, instanceDir := range instanceDirs {

		instance, err := instanceRepository.FindById(instanceDir.Name())
		log.Printf("ALL INSTANCES-----instance name: %s", instanceDir.Name())

		if err != nil {
			log.Printf("ALL INSTANCES-----err finding dir name: %s : err: %v", instanceDir.Name(), err)
			return instances, err
		}

		log.Printf("ALL INSTANCES---------- append instance %v", instance)
		instances = append(instances, instance)
	}

	return instances, nil
}

func (instanceRepository *FileSystemInstanceRepository) createConfig(logstashConf map[string]interface{}, templateFile, outputFile string) error {

	data := map[string]interface{}{
		"logstash": logstashConf,
	}

	f, err := os.Create(outputFile)
	defer f.Close()

	tc, err := gerb.ParseFile(true, templateFile)
	if err != nil {
		log.Println("executing template:", err)
		return err
	}
	tc.Render(f, data)

	return nil
}
