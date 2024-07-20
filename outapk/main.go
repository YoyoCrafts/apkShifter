package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"test/common"
)

type conf struct {
	Uuids   map[string]string `yaml:"uuids"`
	Input   string            `yaml:"input"`
	AppName string            `yaml:"appName"`
}

func walle(input string, channel interface{}, apkPathName, apkName string) (err error) {
	outPath, err := filepath.Abs("out")
	if err != nil {
		logrus.Error(err)
		return
	}
	outPath = outPath + "/" + apkPathName
	if !common.PathExists(outPath) {
		os.MkdirAll(outPath, 0755)
	}

	out := outPath + "/" + apkName + ".apk"
	if common.PathExists(out) {
		os.Remove(out)
	}

	walleFile, err := filepath.Abs("../config/library/walle-cli-all.jar")
	args := make([]string, 0)
	args = append(args, "-jar")
	args = append(args, walleFile)
	args = append(args, "put")
	args = append(args, "-c")
	args = append(args, fmt.Sprint(channel))
	args = append(args, input)
	args = append(args, out)

	var output []byte
	cmd := exec.Command("java", args...)
	output, err = cmd.CombinedOutput()
	logrus.Debug(string(output))
	logrus.Info(fmt.Sprintf("java %s", strings.Trim(fmt.Sprint(args), "[]")))
	if err != nil {
		logrus.Error(fmt.Sprintf("java %s", strings.Trim(fmt.Sprint(args), "[]")))
		logrus.Error(string(output))
		return
	}

	args = make([]string, 0)
	args = append(args, "-jar")
	args = append(args, walleFile)
	args = append(args, "show")
	args = append(args, out)

	cmd = exec.Command("java", args...)
	output, err = cmd.CombinedOutput()
	logrus.Info(fmt.Sprintf("java %s output:%s", strings.Trim(fmt.Sprint(args), "[]"), string(output)))

	return
}

func main() {
	instance := &conf{}
	yamlFile, err := os.ReadFile("./config.yaml")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = yaml.Unmarshal(yamlFile, &instance)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for k, v := range instance.Uuids {
		err = walle(instance.Input, v, k, instance.AppName)
		if err != nil {
			logrus.Error("渠道包绑定ID失败  name:" + k + " channelName:" + v)
			logrus.Error(err)
			return
		}
	}

}
