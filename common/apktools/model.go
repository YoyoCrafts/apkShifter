package apktools

import (
	"crypto/md5"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"sync"
	"test/common"
	"test/common/cache"
	"time"
)

type ReplacePackage struct {
	OldPackageNameKey string
	NewPackageNameKey string
}

var once sync.Once
var Cache *ReplacePackage
var lock sync.Mutex

func ReplacePackageData() *ReplacePackage {
	once.Do(func() {
		Cache = &ReplacePackage{}
	})
	return Cache
}

func TimerDelFile(file string, Minute int) {
	go func() {
		time.Sleep(time.Second * time.Duration(Minute))
		DelFile(file)
		logrus.Debug("删除临时包 file:" + file)
	}()
}

func (c *ReplacePackage) PackageIng() {

	var md5Code = md5.Sum([]byte(common.GetConf().App.ApkPath))
	c.OldPackageNameKey = fmt.Sprintf("OldPackageNameKey:%x", md5Code)
	c.NewPackageNameKey = fmt.Sprintf("NewPackageNameKey:%x", md5Code)

	newPackageName := fmt.Sprintf("com.%s.%s.%s%d%d%d", RandomString(6), RandomString(6), RandomString(3), time.Now().Month(), time.Now().Day(), time.Now().Hour())

	var err error
	var NewPackageNamePath string
	if common.GetConf().App.UpdateConfig.ReplacePackageName {
		NewPackageNamePath, err = SetPackageName(common.GetConf().App.ApkPath, newPackageName)
		if err != nil {
			logrus.Error("更换apk 包名失败 file:" + common.GetConf().App.ApkPath)
			logrus.Error(err)
			return
		}
	} else {
		NewPackageNamePath, err = PackageSign(common.GetConf().App.ApkPath, false)
		if err != nil {
			logrus.Error("更换apk 签名失败 file:" + common.GetConf().App.ApkPath)
			logrus.Error(err)
			return
		}
	}

	OldPackageNamePath := cache.Get(c.OldPackageNameKey)
	if len(OldPackageNamePath) != 0 {
		TimerDelFile(OldPackageNamePath, 120)
		TimerDelFile(OldPackageNamePath+".idsig", 120)
	}

	OldPackageNamePath = cache.Get(c.NewPackageNameKey)
	if len(OldPackageNamePath) != 0 {
		cache.Set(c.OldPackageNameKey, OldPackageNamePath, 60*60*24)
	}

	cache.Set(c.NewPackageNameKey, NewPackageNamePath, 60*60*24)

}

func (c *ReplacePackage) GetChannelApkPath(channelName string) (filePath string, err error) {

	NewPackagePath := cache.Get(c.NewPackageNameKey)
	OldPackagePath := cache.Get(c.OldPackageNameKey)

	md5Code := "sourcePackage"
	if len(NewPackagePath) > 0 {
		md5Code = fmt.Sprintf("%x", md5.Sum([]byte(NewPackagePath)))
	} else if len(OldPackagePath) > 0 {
		md5Code = fmt.Sprintf("%x", md5.Sum([]byte(OldPackagePath)))
	}

	cacheKey := fmt.Sprintf("newChannelWalleFile:%s:%s", md5Code, channelName)
	channelWalleFile := cache.Get(cacheKey)
	if len(channelWalleFile) != 0 && common.PathExists(channelWalleFile) {
		filePath = channelWalleFile
		logrus.Info("下载缓存渠道签名包 channelName:" + channelName + " file:" + channelWalleFile + " :cacheKey:" + cacheKey)
		return
	}

	filePath, err = c.WalleStart(channelName, cacheKey, NewPackagePath, OldPackagePath)
	return

}

func (c *ReplacePackage) GetDowApkPath() (filePath string, err error) {

	NewPackagePath := cache.Get(c.NewPackageNameKey)
	OldPackagePath := cache.Get(c.OldPackageNameKey)

	filePath = common.GetConf().App.ApkPath
	if common.PathExists(NewPackagePath) {
		filePath = NewPackagePath
		logrus.Info("使用最新包名 分包  channelName:" + " file:" + filePath)
	} else if common.PathExists(OldPackagePath) {
		filePath = OldPackagePath
		logrus.Info("使用上一个包名 分包  channelName:" + " file:" + filePath)
	} else {
		if !common.PathExists(filePath) {
			err = errors.New("请配置好apk路径")
			return
		}

		logrus.Info("使用未改名包 分包  channelName:" + " file:" + filePath)
	}

	return

}

func (c *ReplacePackage) WalleStart(channelName string, cacheKey, newPackagePath, oldPackagePath string) (filePath string, err error) {

	channelWalleFile := cache.Get(cacheKey)

	if len(channelWalleFile) != 0 && common.PathExists(channelWalleFile) {
		filePath = channelWalleFile
		return
	}

	packageNameFile := common.GetConf().App.ApkPath
	if common.PathExists(newPackagePath) {
		packageNameFile = newPackagePath
		logrus.Info("使用最新包名 分包  channelName:" + channelName + " file:" + packageNameFile)
	} else if common.PathExists(oldPackagePath) {
		packageNameFile = oldPackagePath
		logrus.Info("使用上一个包名 分包  channelName:" + channelName + " file:" + packageNameFile)
	} else {

		if !common.PathExists(packageNameFile) {
			err = errors.New("请配置好apk路径")
			return
		}

		logrus.Info("使用未改名包 分包  channelName:" + channelName + " file:" + packageNameFile)
	}

	filePath, err = Walle(packageNameFile, channelName)
	if err != nil {
		logrus.Error("渠道包绑定ID失败  channelName:" + channelName + " file:" + filePath)
		logrus.Error(err)
		return
	}

	if !common.PathExists(filePath) {
		err = errors.New("没有找到 渠道包 channelName:" + channelName + " file:" + filePath)
		return
	}

	cache.Set(cacheKey, filePath, 60*30)
	packagenametime := common.GetConf().App.UpdateConfig.Interval
	defer TimerDelFile(filePath, packagenametime+(60*2))
	return
}
