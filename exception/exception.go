package exception

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func Try(userFn func(), catchFn func(err interface{})) {
	defer func() {
		if err := recover();err != nil{
			logrus.Warn(string(debug.Stack()))
			if catchFn != nil {
				catchFn(err)
			}else{
				logrus.Error(fmt.Sprintf("程序执行严重错误: %v", err))
			}
		}
	}()

	userFn()
}

func TryFn(userFn func()) {
	defer func() {
		if err := recover();err != nil{
			logrus.Error(string(debug.Stack()))
			logrus.Error(fmt.Sprintf("程序执行严重错误: %v", err))
		}
	}()
	userFn()
}