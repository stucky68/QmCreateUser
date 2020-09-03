package Utils

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"unicode"
)

var info *log.Logger

func init() {
	file, err := os.OpenFile("./log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}
	info = log.New(io.MultiWriter(file, os.Stderr), "INFO: ", log.Ldate|log.Ltime)
}

func Log(v ...interface{}) {
	info.Println(v)
}

func FilterNickName(content string) string {
	newContent := ""
	for _, value := range content {
		if unicode.Is(unicode.Han, value) || (value >= 'a' && value <= 'z') || (value >= 'A' && value <= 'Z') || unicode.IsDigit(value) || unicode.IsSpace(value) {
			newContent += string(value)
		}
	}
	return newContent
}

func ReadFileData(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}

