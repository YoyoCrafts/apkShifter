package apktools

import (
	"crypto/md5"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"test/common"
	"time"
)

type KeyStoreInfo struct {
	KeyStorePath      string
	KeyStorePass      string
	KeyStoreAlias     string
	KeyStoreAliasPass string
}

// RandomString 生成长度为 length 的随机字符串
func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	letters := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func DelFile(filepath string) {
	if common.PathExists(filepath) {
		os.Remove(filepath)
	}
}

func Walle(input string, out string, channel interface{}) (err error) {
	DelFile(out)
	walle, err := filepath.Abs("config/library/walle-cli-all.jar")
	args := make([]string, 0)
	args = append(args, "-jar")
	args = append(args, walle)
	args = append(args, "put")
	args = append(args, "-c")
	args = append(args, fmt.Sprint(channel))
	args = append(args, input)
	args = append(args, out)

	var output []byte
	cmd := exec.Command("java", args...)
	output, err = cmd.CombinedOutput()
	logrus.Debug(string(output))
	logrus.Info(fmt.Sprintf("java %s", strings.Trim(fmt.Sprint(args), "[]")))
	if err != nil {
		logrus.Error(fmt.Sprintf("java %s", strings.Trim(fmt.Sprint(args), "[]")))
		logrus.Error(string(output))
		return
	}
	return
}

func Zipalign(input string, out string) (err error) {
	DelFile(out)

	sysType := runtime.GOOS

	zipalign, err := filepath.Abs("config/library/zipalign")
	if sysType == "darwin" {
		zipalign, err = filepath.Abs("config/library/zipalign_mac")
	}

	args := make([]string, 0)
	args = append(args, "-v")
	args = append(args, "4")
	args = append(args, input)
	args = append(args, out)

	var output []byte
	cmd := exec.Command(zipalign, args...)
	output, err = cmd.CombinedOutput()
	logrus.Debug(string(output))
	logrus.Info(fmt.Sprintf("%s %s", zipalign, strings.Trim(fmt.Sprint(args), "[]")))

	if err != nil {
		logrus.Error(fmt.Sprintf("%s %s", zipalign, strings.Trim(fmt.Sprint(args), "[]")))
		logrus.Error(string(output))
		return
	}

	args = make([]string, 0)
	args = append(args, "-c")
	args = append(args, "-v")
	args = append(args, "4")
	args = append(args, out)

	cmd = exec.Command(zipalign, args...)
	output, err = cmd.CombinedOutput()
	outputString := string(output)
	logrus.Debug(outputString)
	logrus.Info(fmt.Sprintf("%s %s", zipalign, strings.Trim(fmt.Sprint(args), "[]")))
	if err != nil {
		logrus.Error(fmt.Sprintf("%s %s", zipalign, strings.Trim(fmt.Sprint(args), "[]")))
		logrus.Error(outputString)
		return
	}
	if !strings.Contains(outputString,"Verification succesful"){
		err = fmt.Errorf("zipalign 验证不通过")
	}

	return
}

func TestDecompile(input ,packageName,tempPath string) (out string,err error){
	var apktool string
	apktool, err = filepath.Abs("config/library/apktool")
	out = tempPath + "/" + packageName + "_apk"
	if common.PathExists(out) {
		cmd := exec.Command("rm", "-rf", out)
		cmd.CombinedOutput()
	}

	if err != nil {
		return
	}
	args := make([]string, 0)
	args = append(args, "d")
	args = append(args, input)
	args = append(args, "-o")
	args = append(args, out)
	args = append(args, "--only-main-classes")

	var output []byte
	cmd := exec.Command(apktool, args...)
	output, err = cmd.CombinedOutput()
	logrus.Debug(string(output))
	logrus.Info(fmt.Sprintf("%s %s", apktool, strings.Trim(fmt.Sprint(args), "[]")))
	if err != nil {
		logrus.Error(fmt.Sprintf("%s %s", apktool, strings.Trim(fmt.Sprint(args), "[]")))
		logrus.Error(string(output))
		return
	}
	return
}

func TestDecompilePack(input,packageName,tempPath string)(out string,err error){

	var apktool string
	apktool, err = filepath.Abs("config/library/apktool")
	bao := tempPath + "/" + packageName + "_bao.apk"
	if common.PathExists(bao) {
		DelFile(bao)
	}


	args := make([]string, 0)
	args = append(args, "b")
	args = append(args, "-o")
	args = append(args, bao)
	args = append(args, input)

	var output []byte
	cmd := exec.Command(apktool, args...)
	output, err = cmd.CombinedOutput()
	logrus.Debug(string(output))
	logrus.Info(fmt.Sprintf("%s %s", apktool, strings.Trim(fmt.Sprint(args), "[]")))
	if err != nil {
		logrus.Error(fmt.Sprintf("%s %s", apktool, strings.Trim(fmt.Sprint(args), "[]")))
		logrus.Error(string(output))
		return
	}
	defer os.Remove(bao)


	keyStorePath := tempPath + "/keystore"
	if !common.PathExists(keyStorePath) {
		os.MkdirAll(keyStorePath, 0755)
	}
	keyStoreFile := keyStorePath + "/keystore_" + packageName + ".jks"

	keyStoreInfo, err := CreateKeyStore(keyStoreFile)
	if err != nil {
		return
	}
	defer os.Remove(keyStoreFile)

	zipFile := tempPath + "/" + packageName + "_zip.apk"
	if common.PathExists(zipFile) {
		DelFile(zipFile)
	}
	err = Zipalign(bao, zipFile)
	if err != nil {
		return
	}

	defer os.Remove(zipFile)

	out = tempPath + "/" + packageName + "_sign.apk"
	if common.PathExists(out) {
		DelFile(out)
	}
	err = StartSigning(zipFile, out, keyStoreInfo)
	return
}

func TestStartSigning(inputPath string, outPath string,tempPath string)(err error){

	keyStorePath := tempPath + "/keystore"
	if !common.PathExists(keyStorePath) {
		os.MkdirAll(keyStorePath, 0755)
	}
	keyStoreFile := keyStorePath + "/keystore_TestStartSigning.jks"

	keyStoreInfo, err := CreateKeyStore(keyStoreFile)
	if err != nil {
		return
	}
	zipFile := tempPath + "/TestStartSigning_zip.apk"
	if common.PathExists(zipFile) {
		DelFile(zipFile)
	}
	err = Zipalign(inputPath, zipFile)
	if err != nil {
		return
	}

	err = StartSigning(zipFile, outPath, keyStoreInfo)
	return
}

func SetPackageName(input string, newPackageName string) (signFile string, err error) {
	tempPath, err := filepath.Abs("temp")
	if err != nil {
		return
	}
	tempPath = tempPath + "/source"
	if !common.PathExists(tempPath) {
		os.MkdirAll(tempPath, 0755)
	}

	md5Code := md5.Sum([]byte(input))
	pathKey := fmt.Sprintf("%x", md5Code)

	compilePath := tempPath + "/" + pathKey + "_apk"
	out := tempPath + "/" + pathKey + "_bao.apk"

	var apktool string
	apktool, err = filepath.Abs("config/library/apktool")

	if common.PathExists(compilePath) && common.PathExists(out) {
		DelFile(out)
	} else {
		if common.PathExists(compilePath) {
			cmd := exec.Command("rm", "-rf", compilePath)
			cmd.CombinedOutput()
		}
		if common.PathExists(out) {
			DelFile(out)
		}

		if err != nil {
			return
		}
		args := make([]string, 0)
		args = append(args, "d")
		args = append(args, input)
		args = append(args, "-o")
		args = append(args, compilePath)
		args = append(args, "--only-main-classes")

		var output []byte
		cmd := exec.Command(apktool, args...)
		output, err = cmd.CombinedOutput()
		logrus.Debug(string(output))
		logrus.Info(fmt.Sprintf("%s %s", apktool, strings.Trim(fmt.Sprint(args), "[]")))
		if err != nil {
			logrus.Error(fmt.Sprintf("%s %s", apktool, strings.Trim(fmt.Sprint(args), "[]")))
			logrus.Error(string(output))
			return
		}
	}

	old_package_name, err := common.FileFindAllS(compilePath+"/AndroidManifest.xml", `package="(.*?)"`)
	if err != nil {
		return
	}
	if len(old_package_name) == 0 {
		err = errors.New("没有找到AndroidManifest.xml中的包名")
		return
	}
	//old_package_name_tow := fmt.Sprintf("package=\"%s\"",old_package_name)
	//newPackageName_tow   := fmt.Sprintf("package=\"%s\"",newPackageName)
	//err = common.ReplaceFileContents(compilePath + "/AndroidManifest.xml",old_package_name_tow,newPackageName_tow)

	err = common.ReplaceFileContents(compilePath+"/AndroidManifest.xml", old_package_name, newPackageName)

	args := make([]string, 0)
	args = append(args, "b")
	args = append(args, "-o")
	args = append(args, out)
	args = append(args, compilePath)

	var output []byte
	cmd := exec.Command(apktool, args...)
	output, err = cmd.CombinedOutput()
	logrus.Debug(string(output))
	logrus.Info(fmt.Sprintf("%s %s", apktool, strings.Trim(fmt.Sprint(args), "[]")))
	if err != nil {
		logrus.Error(fmt.Sprintf("%s %s", apktool, strings.Trim(fmt.Sprint(args), "[]")))
		logrus.Error(string(output))
		return
	}

	keyStorePath := tempPath + "/keystore"
	if !common.PathExists(keyStorePath) {
		os.MkdirAll(keyStorePath, 0755)
	}
	keyStoreFile := keyStorePath + "/keystore_" + newPackageName + uuid.NewV4().String() + ".jks"
	keyStoreInfo, err := CreateKeyStore(keyStoreFile)
	if err != nil {
		return
	}
	defer os.Remove(keyStoreFile)
	zipFile := tempPath + "/" + pathKey + uuid.NewV4().String() + "_zip.apk"
	err = Zipalign(out, zipFile)
	if err != nil {
		return
	}
	defer os.Remove(zipFile)

	signFile = tempPath + "/" + pathKey + newPackageName + "_sign.apk"
	err = StartSigning(zipFile, signFile, keyStoreInfo)
	//err = StartSigning(out,signFile,keyStoreInfo)
	return

}

// 开始签名
func StartSigning(inputPath string, outPath string, keyStoreInfo KeyStoreInfo) (err error) {
	DelFile(outPath)
	DelFile(outPath + ".idsig")

	apksigner, err := filepath.Abs("config/library/apksigner.jar")
	if err != nil {
		return
	}
	args := make([]string, 0)
	args = append(args, "-jar")
	args = append(args, apksigner)
	args = append(args, "sign")
	args = append(args, "-verbose")
	args = append(args, "--ks")
	args = append(args, keyStoreInfo.KeyStorePath)
	args = append(args, "--v1-signing-enabled")
	args = append(args, "true")
	args = append(args, "--v2-signing-enabled")
	args = append(args, "true")
	args = append(args, "--ks-pass")
	args = append(args, "pass:"+keyStoreInfo.KeyStorePass)
	args = append(args, "--ks-key-alias")
	args = append(args, keyStoreInfo.KeyStoreAlias)
	args = append(args, "--key-pass")
	args = append(args, "pass:"+keyStoreInfo.KeyStoreAliasPass)
	args = append(args, "--out")
	args = append(args, outPath)
	args = append(args, inputPath)
	var output []byte
	cmd := exec.Command("java", args...)
	output, err = cmd.CombinedOutput()
	logrus.Debug(string(output))
	logrus.Info(fmt.Sprintf("java %s", strings.Trim(fmt.Sprint(args), "[]")))
	if err != nil {
		logrus.Error(fmt.Sprintf("java %s", strings.Trim(fmt.Sprint(args), "[]")))
		logrus.Error(string(output))
		return
	}

	args = make([]string, 0)
	args = append(args, "-jar")
	args = append(args, apksigner)
	args = append(args, "verify")
	args = append(args, "--verbose")
	args = append(args, "--print-certs")
	args = append(args, outPath)

	cmd = exec.Command("java", args...)
	output, err = cmd.CombinedOutput()
	logrus.Debug(string(output))
	logrus.Info(fmt.Sprintf("java %s", strings.Trim(fmt.Sprint(args), "[]")))
	if err != nil {
		logrus.Error(fmt.Sprintf("java %s", strings.Trim(fmt.Sprint(args), "[]")))
		logrus.Error(string(output))
		return
	}
	return

}

//  创建签名
func CreateKeyStore(KeyStorePath string) (keyStoreInfo KeyStoreInfo, err error) {

	DelFile(KeyStorePath)
	pass := RandomString(8)
	keyStoreInfo = KeyStoreInfo{
		KeyStorePath:      KeyStorePath,
		KeyStorePass:      pass,
		KeyStoreAlias:     RandomString(4),
		KeyStoreAliasPass: pass,
	}

	if err != nil {
		return
	}
	args := make([]string, 0)
	args = append(args, "-genkey")
	args = append(args, "-v")
	args = append(args, "-keyalg")
	args = append(args, "RSA")
	args = append(args, "-keysize")
	args = append(args, "2048")
	args = append(args, "-validity")
	args = append(args, "30000")
	args = append(args, "-dname")
	args = append(args, "cn=client-344-839")
	args = append(args, "-keypass")
	args = append(args, keyStoreInfo.KeyStorePass)
	args = append(args, "-storepass")
	args = append(args, keyStoreInfo.KeyStoreAliasPass)
	args = append(args, "-alias")
	args = append(args, keyStoreInfo.KeyStoreAlias)
	args = append(args, "-keystore")
	args = append(args, keyStoreInfo.KeyStorePath)

	var output []byte
	cmd := exec.Command("keytool", args...)
	output, err = cmd.CombinedOutput()

	logrus.Debug(string(output))
	logrus.Info(fmt.Sprintf("keytool %s", strings.Trim(fmt.Sprint(args), "[]")))

	if err != nil {
		logrus.Error(fmt.Sprintf("keytool %s", strings.Trim(fmt.Sprint(args), "[]")))
		logrus.Error(string(output))
		return
	}
	return

}
