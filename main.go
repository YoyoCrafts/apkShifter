package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"test/common"
	"test/common/apktools"
	"test/exception"
	"time"
)

func timerRun() {

	if !common.GetConf().App.UpdateConfig.IntervalEnable && !common.GetConf().App.UpdateConfig.ReplacePackageNameEnable {
		return
	}
	go func() {
		for {
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

	timerRun()

	install, _ := filepath.Abs("config/library/install.sh")
	cmd := exec.Command("sh", install)
	cmd.CombinedOutput()

	router := gin.Default()

	router.GET("/guide/channel/download/:channelName/*xx", func(c *gin.Context) {

		channelName := c.Param("channelName")

		v, err := url.QueryUnescape(channelName)
		if err == nil {
			channelName = v
		}

		var apktempPath string
		if channelName == "{}" || channelName == "" || strings.ToLower(channelName) == "null" || strings.ToLower(channelName) == "undefined" {
			apktempPath, err = apktools.ReplacePackageData().GetDowApkPath()
			if err != nil {
				c.String(http.StatusBadRequest, err.Error())
			} else {
				c.File(apktempPath)
			}
		} else {
			apktempPath, err = apktools.ReplacePackageData().GetChannelApkPath(channelName)
			if err != nil {
				c.String(http.StatusBadRequest, err.Error())
			} else {
				c.File(apktempPath)
			}
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
