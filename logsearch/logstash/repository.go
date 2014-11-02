package logstash

import (
	"github.com/karlseguin/gerb"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type InstanceRepository interface {
	CreateConfig(templateFile, outputFile string) error
	CreateInstanceDirectories(instance *Instance) error
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

func (instanceRepository *FileSystemInstanceRepository) CreateInstanceDirectories(instance *Instance) error {
	createBaseDirErr := instanceRepository.createBaseDirectory(instance)
	if createBaseDirErr != nil {
		return createBaseDirErr
	}

	createLogDirErr := instanceRepository.createLogDirectory(instance)
	if createLogDirErr != nil {
		return createLogDirErr
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

func (instanceRepository *FileSystemInstanceRepository) instanceDataDirectory() string {
	return instanceRepository.LogstashConf.InstanceDataDirectory
}

func (instanceRepository *FileSystemInstanceRepository) instanceLogDirectory() string {
	return instanceRepository.LogstashConf.InstanceLogDirectory
}

func (instanceRepository *FileSystemInstanceRepository) FindById(instanceId string) (*Instance, error) {
	instanceDataDir := path.Join(instanceRepository.instanceDataDirectory(), instanceId)

	_, err := os.Stat(instanceDataDir)
	if err != nil {
		return nil, err
	}

	instance := &Instance{
		Id:   instanceId,
		Port: 5514,
		Host: instanceRepository.LogstashConf.Host,
	}

	return instance, nil
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

func (instanceRepository *FileSystemInstanceRepository) CreateConfig(templateFile, outputFile string) error {

	data := map[string]interface{}{
		"logstash": map[string]interface{}{"Host": "127.0.0.1", "Port": 5514},
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

func (instanceRepository *FileSystemInstanceRepository) GetInstanceCount() (int, error) {
	instances, err := instanceRepository.findAllInstances()
	return len(instances), err
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
