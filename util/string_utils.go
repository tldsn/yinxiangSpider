package util

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//URLEncode ...
func URLEncode(str string) string {
	str = url.QueryEscape(str)
	str = strings.Replace(str, "%3D", "=", -1)
	str = strings.Replace(str, "%26", "&", -1)
	return str
}

//SaveFile ...
func SaveFile(path, content string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(content))
	if err != nil {
		return err
	}
	return nil
}

//GetLogger ...
func GetLogger(path string) (logger *log.Logger) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 666)
	if err != nil {
		log.Fatalln("fail to create test.log file!")
	}
	// defer file.Close()
	logger = log.New(file, "", log.LstdFlags|log.Lshortfile) // 日志文件格式:log包含时间及文件行数
	logger.SetFlags(log.LstdFlags | log.Lshortfile)          // 设置日志格式
	return logger
}

//BetweenString ...
func BetweenString(str, beg, end string) string {
	index := strings.Index(str, beg)
	if index == -1 {
		return ""
	}
	tmp := str[index+len(beg):]
	index = strings.Index(tmp, end)
	if index == -1 {
		return ""
	}
	return tmp[:index]
}

//RandomUUID ...
func RandomUUID() string {
	var u [16]byte
	if _, err := io.ReadFull(rand.Reader, u[:]); err != nil {
		return ""
	}
	log.Println(u[:])
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}

//GetCurrentDir 获取当前执行程序所在路径
func GetCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	if runtime.GOOS == "windows" {
		return dir + "\\"
	}
	return dir + "/"
}
