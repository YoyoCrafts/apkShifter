package common

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

func PathExists(path string) bool {
	if len(path) == 0 {
		return false
	}
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func FileFindAllS(fileName string, str string) (re string, err error) {

	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	pos := int64(0)
	for {
		var line string
		line, err = reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			} else {
				return
			}
		}

		res := regexp.MustCompile(str)
		alls := res.FindAllStringSubmatch(line, -1)
		if len(alls) > 0 {
			re = alls[0][1]
			return
		}

		pos += int64(len(line))
	}

	return
}

func ReplaceFileContents(fileName string, oldString string, newString string) (err error) {

	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	out, err := os.OpenFile(fileName+".mdf", os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		fmt.Println("Open write file fail:", err)
		os.Exit(-1)
	}
	defer out.Close()

	br := bufio.NewReader(file)
	index := 1
	for {
		var line []byte
		line, _, err = br.ReadLine()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			os.Exit(-1)
			return
		}
		newLine := strings.ReplaceAll(string(line), oldString, newString)
		_, err = out.WriteString(newLine + "\n")
		if err != nil {
			os.Exit(-1)
			return
		}
		index++
	}

	DelFile(fileName)
	if err == nil {
		err = os.Rename(fileName+".mdf", fileName)
	}

	defer DelFile(fileName + ".mdf")

	return
}

// RandomString 生成长度为 length 的随机字符串
func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

func DelFile(filepath string) {
	if PathExists(filepath) {
		os.Remove(filepath)
	}
}
