package common

import (
	"bytes"
	"fmt"
	"github.com/gookit/color"
	 rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"time"
)


type MyFormatter struct{
	colour bool
}



func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string

	if entry.HasCaller() {
		fName := entry.Caller.File
		newLog = fmt.Sprintf("[%s] [%s] [%s:%d] [%s]\n\n",timestamp,entry.Level, fName, entry.Caller.Line, entry.Message)
	} else {
		newLog = fmt.Sprintf("[%s] [%s] [%s]\n\n",timestamp,entry.Level,entry.Message)
	}
	//newLog = entry.Message

	if m.colour && (entry.Level == logrus.ErrorLevel || entry.Level == logrus.PanicLevel || entry.Level == logrus.FatalLevel ) {
		newLog = red(newLog)
	}
	if m.colour && entry.Level == logrus.WarnLevel {
		newLog = yellow(newLog)
	}
	if m.colour && entry.Level == logrus.InfoLevel {
		newLog = green(newLog)
	}

	if m.colour && entry.Level == logrus.DebugLevel {
		newLog = blue(newLog)
	}


	b.WriteString(newLog)
	return b.Bytes(), nil
}

func newRotatelogs(fileName string) *rotatelogs.RotateLogs {
	logier, err := rotatelogs.New(
		fileName,
		rotatelogs.WithMaxAge(5*24*time.Hour),    // 文件最大保存时间
		rotatelogs.WithRotationTime(1*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		panic(err)
	}
	return logier
}

// 配置日志切割
func LogFileCut(logPath string) {

	errPaht     := logPath+"/error"
	otherPath   := logPath+"/other"
	warnPath   := logPath+"/warn"

	if !PathExists(logPath){
		os.MkdirAll(logPath, 0755)
	}
	if !PathExists(errPaht){
		os.MkdirAll(errPaht, 0755)
	}
	if !PathExists(otherPath){
		os.MkdirAll(otherPath, 0755)
	}
	if !PathExists(warnPath){
		os.MkdirAll(warnPath, 0755)
	}

	otherFileName := path.Join(otherPath, "%Y-%m-%d .log")
	errFileName   := path.Join(errPaht,   "%Y-%m-%d .log")
	warnFileName   := path.Join(warnPath, "%Y-%m-%d .log")


	errlogier   := newRotatelogs(errFileName)
	otherlogier := newRotatelogs(otherFileName)
	warnlogier := newRotatelogs(warnFileName)


	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel:  otherlogier,
		logrus.FatalLevel: otherlogier,
		logrus.DebugLevel: otherlogier,
		logrus.WarnLevel:  warnlogier,
		logrus.ErrorLevel: errlogier,
		logrus.PanicLevel: errlogier,
	},&MyFormatter{})
	logrus.AddHook(lfHook)


	writers := []io.Writer{ errlogier, os.Stdout,otherlogier, os.Stdout}
	io.MultiWriter(writers...)

}

func InitLog() {
	logInfo := GetConf().Log

	logrus.SetReportCaller(true)
	logrus.SetFormatter(&MyFormatter{
		colour:logInfo.Colour,
	})

	level,_ := logrus.ParseLevel(logInfo.Level)
	logrus.SetLevel(level)

	// 配置日志分割
	LogFileCut(logInfo.File)

}




func red(msg string)string  {
	return color.Red.Sprintf(msg)
}

func yellow(msg string)string  {
	return color.Yellow.Sprintf(msg)
}

func green(msg string)string  {
	return color.Green.Sprintf(msg)
}

func blue(msg string)string  {
	return color.Blue.Sprintf(msg)
}

