package main

import (
	"github.com/sirupsen/logrus"
	"test/common/apktools"
)

func main() {


	//_,err := apktools.TestDecompile("/Users/xiaobai/Downloads/7000/双端/A命运女神.apk","video.test.xiaobai.com","/Users/xiaobai/Downloads/7000/双端/temp")
	//if err != nil{
	//	logrus.Error(err)
	//	return
	//}
	_,err := apktools.TestDecompilePack("/Users/xiaobai/Downloads/7000/双端/temp/video.test.xiaobai.com_apk","video.test.xiaobai.com","/Users/xiaobai/Downloads/7000/双端/temp")
	if err != nil{
		logrus.Error(err)
		return
	}


	//err := apktools.TestStartSigning("/Users/xiaobai/Downloads/7000/双端/A命运女神2.apk","/Users/xiaobai/Downloads/7000/双端/A命运女神2_qm.apk","/Users/xiaobai/Downloads/7000/双端/temp")
	//if err != nil{
	//	logrus.Error(err)
	//	return
	//}



	//err := ListDir("/Users/xiaobai/Downloads/7000/yuanm/ym/client的副本")
	//if err != nil{
	//	logrus.Error(err)
	//	return
	//}

}

//
//func ListDir(dirname string) (error) {
//	infos, err := ioutil.ReadDir(dirname)
//	if err != nil {
//		return err
//	}
//	for _, info := range infos {
//		if info.IsDir() {
//			err = ListDir(dirname+"/"+info.Name() )
//			if err != nil{
//				return err
//			}
//		}else if path.Ext(info.Name()) == ".luac"{
//			ext := path.Ext(info.Name())
//			os.Rename(dirname+"/"+info.Name(), dirname+"/"+strings.Replace(info.Name(),ext,".lua",-1))
//		}
//	}
//	return  nil
//}