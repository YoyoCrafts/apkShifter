package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"test/common"
	"test/common/apktools"
	"test/exception"
	"time"
)

func timerRun() {

	go func() {
		for {

			if !common.GetConf().App.UpdateConfig.Switch {
				break
			}
			exception.TryFn(func() {
				apktools.ReplacePackageData().PackageIng()
			})
			time.Sleep(time.Duration(common.GetConf().App.UpdateConfig.Interval) * time.Second)
		}
	}()
}

func main() {

	if !common.PathExists("temp/cache") {
		os.MkdirAll("temp/cache", 0755)
	}
	common.InitLog()
	//apktools.StartSigningVerify(common.GetConf().App.ApkPath)
	timerRun()

	install, _ := filepath.Abs("config/library/install.sh")
	cmd := exec.Command("sh", install)
	cmd.CombinedOutput()

	router := gin.Default()

	router.GET("/guide/channel/download/:channelName/*xx", func(c *gin.Context) {

		channelName := c.Param("channelName")

		apktempPath, err := apktools.ReplacePackageData().GetChannelApkPath(channelName)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		} else {
			c.File(apktempPath)
		}
	})

	router.GET("/guide/download/*xx", func(c *gin.Context) {

		apktempPath, err := apktools.ReplacePackageData().GetDowApkPath()
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		} else {
			c.File(apktempPath)
		}
	})

	router.Run(":" + strconv.Itoa(common.GetConf().App.Port))

}

//GOOS=linux GOARCH=amd64 go build main.go
