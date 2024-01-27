package apktools

import (
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

func Walle(input string, channel interface{}) (out string, err error) {
	tempPath, err := filepath.Abs("temp")
	if err != nil {
		logrus.Error(err)
		return
	}
	tempPath = tempPath + "/walle"
	if !common.PathExists(tempPath) {
		os.MkdirAll(tempPath, 0755)
	}

	out = tempPath + "/" + uuid.NewV4().String() + fmt.Sprint(channel) + ".apk"
	if common.PathExists(out) {
		os.Remove(out)
	}

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

func Zipalign(input string) (out string, err error) {

	tempPath, err := filepath.Abs("temp")
	if err != nil {
		logrus.Error(err)
		return
	}

	tempPath = tempPath + "/zipalign"
	if !common.PathExists(tempPath) {
		os.MkdirAll(tempPath, 0755)
	}

	out = tempPath + "/" + uuid.NewV4().String() + "zip.apk"
	if common.PathExists(out) {
		os.Remove(out)
	}
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
	if !strings.Contains(outputString, "Verification succesful") {
		err = fmt.Errorf("zipalign 验证不通过")
	}

	return
}

// 替换apk包名并且签名
func SetPackageName(input string, newPackageName string) (signFile string, err error) {
	tempPath, err := filepath.Abs("temp")
	if err != nil {
		return
	}
	tempPath = tempPath + "/source/" + uuid.NewV4().String()
	if !common.PathExists(tempPath) {
		os.MkdirAll(tempPath, 0755)
	}

	compilePath := tempPath + "/compilePath"

	out := tempPath + "/newPackageName.apk"

	var apktool string
	apktool, err = filepath.Abs("config/library/apktool")

	if common.PathExists(tempPath) {
		cmd := exec.Command("rm", "-rf", tempPath)
		cmd.CombinedOutput()
	}
	if common.PathExists(out) {
		os.Remove(out)
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

	old_package_name, err := common.FileFindAllS(compilePath+"/AndroidManifest.xml", `package="(.*?)"`)
	if err != nil {
		return
	}
	if len(old_package_name) == 0 {
		err = errors.New("没有找到AndroidManifest.xml中的包名")
		return
	}
	err = common.ReplaceFileContents(compilePath+"/AndroidManifest.xml", old_package_name, newPackageName)

	args = make([]string, 0)
	args = append(args, "b")
	args = append(args, "-o")
	args = append(args, out)
	args = append(args, compilePath)

	cmd = exec.Command(apktool, args...)
	output, err = cmd.CombinedOutput()
	logrus.Debug(string(output))
	logrus.Info(fmt.Sprintf("%s %s", apktool, strings.Trim(fmt.Sprint(args), "[]")))
	if err != nil {
		logrus.Error(fmt.Sprintf("%s %s", apktool, strings.Trim(fmt.Sprint(args), "[]")))
		logrus.Error(string(output))
		return
	}

	return PackageSign(out, true)

}

func PackageSign(input string, zip bool) (signFile string, err error) {

	if zip {
		input, err = Zipalign(input)
		if err != nil {
			return
		}
		defer os.Remove(input)
	}

	keyStoreInfo, err := CreateKeyStore()
	if err != nil {
		return
	}

	signFile, err = StartSigning(input, keyStoreInfo)

	defer os.Remove(keyStoreInfo.KeyStorePath)
	return
}

// 开始签名
func StartSigning(inputPath string, keyStoreInfo KeyStoreInfo) (outPath string, err error) {

	tempPath, err := filepath.Abs("temp")
	if err != nil {
		logrus.Error(err)
		return
	}

	tempPath = tempPath + "/sign"
	if !common.PathExists(tempPath) {
		os.MkdirAll(tempPath, 0755)
	}

	outPath = tempPath + "/" + uuid.NewV4().String() + "sign.apk"
	if common.PathExists(outPath) {
		os.Remove(outPath)
		os.Remove(outPath + ".idsig")
	}
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
	args = append(args, "--v3-signing-enabled")
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
	err = StartSigningVerify(outPath)
	if err != nil {
		return
	}

	return

}

func StartSigningVerify(inputPath string) (err error) {
	apksigner, err := filepath.Abs("config/library/apksigner.jar")
	args := make([]string, 0)
	args = make([]string, 0)
	args = append(args, "-jar")
	args = append(args, apksigner)
	args = append(args, "verify")
	args = append(args, "--verbose")
	args = append(args, "--print-certs")
	args = append(args, inputPath)

	cmd := exec.Command("java", args...)
	output, err := cmd.CombinedOutput()
	outputString := string(output)
	logrus.Debug(outputString)
	logrus.Info(fmt.Sprintf("java %s", strings.Trim(fmt.Sprint(args), "[]")))
	if err != nil {
		logrus.Error(fmt.Sprintf("java %s", strings.Trim(fmt.Sprint(args), "[]")))
		logrus.Error(string(output))
		return
	}
	if strings.Contains(outputString, "Verified using v2 scheme (APK Signature Scheme v2): true") &&
		strings.Contains(outputString, "Number of signers: 1") {
		return
	}
	err = fmt.Errorf("签名 验证不通过 \n%s", outputString)
	return
}

// 创建签名
func CreateKeyStore() (keyStoreInfo KeyStoreInfo, err error) {
	tempPath, err := filepath.Abs("temp")
	if err != nil {
		logrus.Error(err)
		return
	}
	tempPath = tempPath + "/keystore"
	if !common.PathExists(tempPath) {
		os.MkdirAll(tempPath, 0755)
	}

	keyStoreFile := tempPath + "/" + uuid.NewV4().String() + "keystore.jks"
	if common.PathExists(keyStoreFile) {
		os.Remove(keyStoreFile)
	}

	pass := RandomString(8)
	keyStoreInfo = KeyStoreInfo{
		KeyStorePath:      keyStoreFile,
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
	args = append(args, keyStoreInfo.KeyStoreAliasPass)
	args = append(args, "-storepass")
	args = append(args, keyStoreInfo.KeyStorePass)
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
