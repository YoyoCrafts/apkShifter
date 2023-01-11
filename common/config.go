package common

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type conf struct {
	App  app
	Log  log
}

type app struct {
	Port     int
	Name     string
	Apkpath  string
	Walleinfo string
	Packagenametime int
}



type log struct {
	File string
	Level string
	Colour bool
	Levelsql string
}



//配置文件配置信息
func GetConf() conf {
	var mConf conf
	yamlFile, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}

	err = yaml.Unmarshal(yamlFile, &mConf)
	if err != nil {
		fmt.Println(err.Error())
	}

	return  mConf
}



