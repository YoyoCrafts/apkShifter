package main

import (
	"github.com/sirupsen/logrus"
	"test/common/apktools"
)

func main() {


	//_,err := apktools.TestDecompile("/Users/xiaobai/Downloads/release/simpleDemo-release.apk","video.test.xiaobai.com","/Users/xiaobai/Downloads/test")
	//if err != nil{
	//	logrus.Error(err)
	//	return
	//}
	_,err := apktools.TestDecompilePack("/Users/xiaobai/Downloads/test/video.test.xiaobai.com_apk","video.test.xiaobai.com","/Users/xiaobai/Downloads/test")
	if err != nil{
		logrus.Error(err)
		return
	}

}
