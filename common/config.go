package common

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"sync"
)

type conf struct {
	App app `yaml:"app"`
	Log log `yaml:"log"`
}

type app struct {
	Port               int    `yaml:"port"`
	ApkPath            string `yaml:"apkPath"`
	UpdateInterval     int    `yaml:"updateInterval"`
	ReplacePackageName bool   `yaml:"replacePackageName"`
}

type log struct {
	File     string `yaml:"file"`
	Level    string `yaml:"level"`
	Colour   bool   `yaml:"colour"`
	Levelsql string `yaml:"levelsql"`
}

var instance *conf
var once sync.Once

func GetConf() *conf {
	once.Do(func() {
		instance = &conf{}
		yamlFile, err := os.ReadFile("config/config.yaml")
		if err != nil {
			fmt.Println(err.Error())
		}

		err = yaml.Unmarshal(yamlFile, &instance)
		if err != nil {
			fmt.Println(err.Error())
		}
		instance.App.ApkPath, _ = filepath.Abs(instance.App.ApkPath)
	})
	return instance
}
