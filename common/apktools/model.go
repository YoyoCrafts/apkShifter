package apktools

import (
	"crypto/md5"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
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

func (c *ReplacePackage) ReplacePackageName() {

	var md5Code = md5.Sum([]byte(common.GetConf().App.Apkpath))
	c.OldPackageNameKey = fmt.Sprintf("OldPackageNameKey:%x", md5Code)
	c.NewPackageNameKey = fmt.Sprintf("NewPackageNameKey:%x", md5Code)

	newPackageName := fmt.Sprintf("com.%s.%s.%s%d%d%d", RandomString(6), RandomString(6), RandomString(3), time.Now().Month(), time.Now().Day(), time.Now().Hour())
	SetPackageNameFile, err := SetPackageName(common.GetConf().App.Apkpath, newPackageName)
	if err != nil {
		logrus.Error("更换apk 包名失败 file:" + common.GetConf().App.Apkpath)
		logrus.Error(err)
		return
	}

	OldPackageNameKey := cache.Get(c.OldPackageNameKey)
	if len(OldPackageNameKey) != 0 {
		TimerDelFile(OldPackageNameKey, 120)
		TimerDelFile(OldPackageNameKey+".idsig", 120)
	}

	NewPackageNameKey := cache.Get(c.NewPackageNameKey)
	if len(NewPackageNameKey) != 0 {
		cache.Set(c.OldPackageNameKey, NewPackageNameKey, 60*60*24)
	}

	cache.Set(c.NewPackageNameKey, SetPackageNameFile, 60*60*24)

	logrus.Info("APK 包名更换成功 新包名:" + newPackageName)
}

func (c *ReplacePackage) GetChannelPath(channelName string) (filePath string, err error) {


	NewPackageNameKey := cache.Get(c.NewPackageNameKey)
	OldPackageNameKey := cache.Get(c.OldPackageNameKey)

	md5Code := ""
	if len(NewPackageNameKey) > 0 {
		md5Code = fmt.Sprintf("%x", md5.Sum([]byte(NewPackageNameKey)))
	}
	cacheKey := fmt.Sprintf("newChannelWalleFile:%s:%s", md5Code, channelName)

	channelWalleFile := cache.Get(cacheKey)
	if len(channelWalleFile) != 0 && common.PathExists(channelWalleFile) {
		filePath = channelWalleFile
		logrus.Info("下载缓存渠道签名包 channelName:" + channelName + " file:" + channelWalleFile + " :cacheKey:" + cacheKey)
		return
	}




	if len(OldPackageNameKey) > 0 {
		md5Code = fmt.Sprintf("%x", md5.Sum([]byte(OldPackageNameKey)))

		cacheKey = fmt.Sprintf("newChannelWalleFile:%s:%s", md5Code, channelName)
		channelWalleFile = cache.Get(cacheKey)
		if len(channelWalleFile) != 0 && common.PathExists(channelWalleFile) {
			filePath = channelWalleFile
			logrus.Info("下载上一个缓存渠道签名包 channelName:" + channelName + " file:" + channelWalleFile + " :cacheKey:" + cacheKey)
			go c.WalleStart(channelName)
			return
		}
	}



	filePath,err = c.WalleStart(channelName)
	return

}

func (c *ReplacePackage)WalleStart(channelName string)  (filePath string, err error){



	NewPackageNameKey := cache.Get(c.NewPackageNameKey)
	OldPackageNameKey := cache.Get(c.OldPackageNameKey)

	lock.Lock()
	defer lock.Unlock()


	md5Code := ""
	if len(NewPackageNameKey) > 0 {
		md5Code = fmt.Sprintf("%x", md5.Sum([]byte(NewPackageNameKey)))
	}
	cacheKey := fmt.Sprintf("newChannelWalleFile:%s:%s", md5Code, channelName)

	channelWalleFile := cache.Get(cacheKey)

	if len(channelWalleFile) != 0 && common.PathExists(channelWalleFile) {
		filePath = channelWalleFile
		return
	}


	tempPath, err := filepath.Abs("temp")
	if err != nil {
		logrus.Error(err)
		return
	}


	wallePath := tempPath + "/wallepath"
	if !common.PathExists(wallePath) {
		os.MkdirAll(wallePath, 0755)
	}

	for {
		filePath = wallePath + "/" + channelName + "_" + uuid.NewV4().String() + ".apk"
		if !common.PathExists(filePath) {
			break
		}
	}

	packageNameFile := common.GetConf().App.Apkpath
	if common.PathExists(NewPackageNameKey) {
		packageNameFile = NewPackageNameKey
		logrus.Info("使用最新包名 分包  channelName:" + channelName + " file:" + packageNameFile)
	} else if common.PathExists(OldPackageNameKey) {
		packageNameFile = OldPackageNameKey
		logrus.Info("使用上一个包名 分包  channelName:" + channelName + " file:" + packageNameFile)
	} else {

		if !common.PathExists(packageNameFile) {
			err = errors.New("请配置好apk路径")
			return
		}

		logrus.Info("使用未改名包 分包  channelName:" + channelName + " file:" + packageNameFile)

		////签名
		keyStorePath := tempPath +"/keystore"
		if !common.PathExists(keyStorePath){
			os.MkdirAll(keyStorePath, 0755)
		}
		keyStoreFile := keyStorePath+"/keystore_"+uuid.NewV4().String()+".jks"

		keyStoreInfo,errs :=  CreateKeyStore(keyStoreFile)
		if errs!=nil {
			logrus.Error("渠道包 生产签名失败  channelName:" + channelName + " file:" + filePath)
			logrus.Error(errs)
			return
		}
		signingFile := strings.TrimRight(filePath, ".apk")+"_sign.apk"


		errs = StartSigning(packageNameFile,signingFile,keyStoreInfo)
		if errs!=nil {
			logrus.Error("渠道包 签名失败  channelName:" + channelName + " file:" + filePath)
			logrus.Error(errs)
			return
		}
		defer TimerDelFile(signingFile, 60)

		packageNameFile = signingFile
	}




	err = Walle(packageNameFile, filePath, channelName)
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
	packagenametime := common.GetConf().App.Packagenametime
	defer TimerDelFile(filePath, packagenametime+(60*2))
	return
}

