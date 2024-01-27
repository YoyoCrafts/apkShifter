package apps

//
//import (
//	"crypto/md5"
//	"encoding/json"
//	"fmt"
//	"github.com/sirupsen/logrus"
//	"image"
//	_ "image/gif"
//	_ "image/jpeg"
//	"image/png"
//	"io"
//	"net/http"
//	"os"
//	"path/filepath"
//	"strings"
//	"test/common"
//	"test/common/apktools"
//	"test/exception"
//	"time"
//)
//
//type AppConfig struct {
//	AppLogo string `json:AppLogo` // path or net
//	AppName string `json:AppName` // app名称
//}
//
//type Config struct {
//	AppConfigs []AppConfig `json:"appConfigs"`
//}
//
//func Copy(input, out string) error {
//	// 打开源文件
//	srcFile, err := os.Open(input)
//	if err != nil {
//		return err
//	}
//	defer srcFile.Close()
//
//	if common.PathExists(out) {
//		os.Remove(out)
//	}
//
//	// 创建目标文件
//	dstFile, err := os.Create(out)
//	if err != nil {
//		return err
//	}
//	defer dstFile.Close()
//
//	// 将源文件内容复制到目标文件
//	_, err = io.Copy(dstFile, srcFile)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func Dow(url, dowfile string) error {
//
//	// 发送 HTTP GET 请求
//	resp, err := http.Get(url)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	if common.PathExists(dowfile) {
//		os.Remove(dowfile)
//	}
//
//	// 解码图片
//	img, _, err := image.Decode(resp.Body)
//	if err != nil {
//		return err
//	}
//
//	// 创建一个新的PNG文件
//	outputFile, err := os.Create(dowfile)
//	if err != nil {
//		return err
//	}
//	defer outputFile.Close()
//
//	// 将图片编码为PNG格式并保存到文件
//	err = png.Encode(outputFile, img)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func Apps() {
//	configFile, err := os.Open("apps/config.json")
//	if err != nil {
//		fmt.Sprintf("请带上配置文件")
//		return
//	}
//
//	defer configFile.Close()
//
//	// 解码 JSON 文件到结构体
//	var configInfo Config
//	err = json.NewDecoder(configFile).Decode(&configInfo)
//	if err != nil {
//		fmt.Sprintf("配置文件解析失败")
//		return
//	}
//
//	if !common.PathExists(common.GetConf().App.Apkpath) {
//		panic("请设置好app的路径")
//	}
//
//	for i, v := range configInfo.AppConfigs {
//
//		apkPath, _ := filepath.Abs(fmt.Sprintf("apps/%s.apk", v.AppName))
//
//		if common.PathExists(apkPath) {
//			continue
//		}
//
//		appTmpPathName := fmt.Sprintf("%x", md5.Sum([]byte(common.GetConf().App.Apkpath)))
//		appTmpPath, _ := filepath.Abs(fmt.Sprintf("temp"))
//		appTmpPath1, _ := filepath.Abs(fmt.Sprintf("%s/%s", appTmpPath, appTmpPathName))
//		if !common.PathExists(appTmpPath1) {
//			err = apktools.TestDecompile(common.GetConf().App.Apkpath, appTmpPath1)
//			if err != nil {
//				logrus.Error(err)
//				return
//			}
//		}
//
//		if strings.Contains(v.AppLogo, "http") {
//			appTmpPathName2 := fmt.Sprintf("%x", md5.Sum([]byte(v.AppLogo)))
//			appTmpPath2, _ := filepath.Abs(fmt.Sprintf("temp/1%s.png", appTmpPathName2))
//			if common.PathExists(appTmpPath2) {
//				os.Remove(appTmpPath2)
//			}
//			err := Dow(v.AppLogo, appTmpPath2)
//			if err != nil {
//				logrus.Errorf("下载图片失败 %s   %s", v.AppLogo, err.Error())
//				continue
//			}
//			v.AppLogo = appTmpPath2
//		}
//
//		logoName := fmt.Sprintf("app_logo_%d", i)
//		drawablePathList := []string{
//			"/res/drawable",
//			"/res/drawable-hdpi",
//			"/res/drawable-ldpi",
//			"/res/drawable-ldrtl-hdpi",
//			"/res/drawable-ldrtl-mdpi",
//			"/res/drawable-ldrtl-xhdpi",
//			"/res/drawable-ldrtl-xxhdpi",
//			"/res/drawable-ldrtl-xxxhdpi",
//			"/res/drawable-mdpi",
//			"/res/drawable-watch",
//			"/res/drawable-xhdpi",
//			"/res/drawable-xxhdpi",
//			"/res/drawable-xxxhdpi",
//		}
//
//		for _, vv := range drawablePathList {
//			logoNamePath := appTmpPath1 + vv + "/" + logoName + ".png"
//			if common.PathExists(appTmpPath1+vv) && !common.PathExists(logoNamePath) {
//				err = Copy(v.AppLogo, logoNamePath)
//				if err != nil {
//					logrus.Errorf("复制logo失败 path:%s err:%s", logoNamePath, err)
//					return
//				}
//			}
//		}
//
//		var signFile string
//		signFile, err = apktools.TestSetAndroidManifest(appTmpPath1, appTmpPathName, v.AppName, logoName)
//		if err != nil {
//			logrus.Error(err)
//			continue
//		}
//
//		if !common.PathExists("apps") {
//			err = os.Mkdir("apps", 0755)
//			if err != nil {
//				return
//			}
//		}
//
//		// 移动文件
//		err = os.Rename(signFile, apkPath)
//		if err != nil {
//			logrus.Error(err)
//			return
//		}
//
//		logrus.Info(apkPath)
//
//	}
//}
//
//
//func AppsRun() {
//	go func() {
//		for {
//			timer := 10
//			exception.TryFn(func() {
//				Apps()
//			})
//			time.Sleep(time.Duration(timer) * time.Second)
//		}
//	}()
//}
