package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
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
			exception.TryFn(func() {
				apktools.ReplacePackageData().ReplacePackageName()
			})
			packagenametime := common.GetConf().App.Packagenametime
			time.Sleep(time.Duration(packagenametime) * time.Second)
		}
	}()
}

func getIpAddr() string {
	addrs, err1 := net.InterfaceAddrs()

	if err1 != nil {
		fmt.Println(err1)
		os.Exit(1)
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {

			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}


func main() {


	if !common.PathExists("temp/cache"){
		os.MkdirAll("temp/cache", 0755)
	}
	common.InitLog()
	timerRun()

	install, _ := filepath.Abs("config/library/install.sh")
	cmd:=exec.Command("sh", install)
	cmd.CombinedOutput()


	router := gin.Default()
	router.GET("/guide/download/:channelName/*xx", func(c *gin.Context) {

		walleinfo := common.GetConf().App.Walleinfo
		if walleinfo == "" {
			walleinfo = "%s"
		}
		channelName := fmt.Sprintf(walleinfo,fmt.Sprint(c.Param("channelName")))

		apktempPath := ""
		apktempPath, err := apktools.ReplacePackageData().GetChannelPath(channelName)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		} else {
			c.File(apktempPath)
		}
	})



	router.Run(":" + strconv.Itoa(common.GetConf().App.Port))


	//for  {
	//	fmt.Printf("APK 防红程序 启动成功 \r\n")
	//	fmt.Printf(fmt.Sprintf("APK 防红下载地址:http://%s:%d/guide/download/channel_id/app_name.apk\n",getIpAddr(),common.GetConf().App.Port))
	//	fmt.Printf("channel_id: walle分包渠道号 如果没有使用walle分包可顺便填写一个字符串\n")
	//	fmt.Printf("app_name:   app名称  随机名称就行\n")
	//
	//	//
	//	//fmt.Printf("您当前使用的是 演示测试版本 10分钟后将会自动退出\n")
	//	//time.Sleep(time.Second * 60 * 10)
	//	//tempPath, _ := filepath.Abs("temp")
	//	//cmd = exec.Command("rm","-rf",tempPath)
	//	//cmd.CombinedOutput()
	//	//fmt.Printf("您当前使用的是 演示测试版本 测试结束 请经快购买正版使用\n")
	//	//break
	//	//
	//
	//	select {}
	//}




}

//env GOOS=linux GOARCH=amd64 go build main.go