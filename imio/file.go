package imio

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func getFileType(url string) string {
	fileIndex := strings.LastIndex(url, ".") + 1
	fileType := url[fileIndex:]
	return fileType
}

func handleFilePost(w http.ResponseWriter, r *http.Request) *AppError {
	if r.Method == "POST" {
		r.ParseForm()
		fileType := "png"
		filePath, fileName := GetFile(fileType)
		f, e := os.Create(filePath)
		if e != nil {
			return &AppError{error: e, statusCode: 500}
		}
		_, err := io.Copy(f, r.Body)
		if err != nil {
			return &AppError{error: err, statusCode: 500}
		}
		url := createUrl(fileName)
		sendOkWithData(w, url)
		defer f.Close()
	} else {
		return &AppError{message: "请求方式错误", statusCode: 400}
	}
	return nil
}

func handleFileGet(w http.ResponseWriter, r *http.Request) *AppError {
	if r.Method == "GET" {
		filePath := getPathFromUrl(r.RequestURI)
		println(filePath)
		file, err := os.Open(filePath)
		if err != nil {
			return &AppError{statusCode: 400, message: "找不到文件", error: err}
		}
		_, err = io.Copy(w, file)
		defer file.Close()
	}
	return nil
}

func GetFile(fileType string) (string, string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	len := 10
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	fileName := string(bytes) + "." + fileType
	basePath := "E:\\file\\type_" + fileType
	_ = os.MkdirAll(basePath, os.ModeAppend)
	filePath := basePath + "\\" + fileName
	return filePath, fileName
}
func createUrl(fileName string) string {
	url := "/file/get/" + fileName
	return url
}

func getPathFromUrl(url string) string {
	index := strings.LastIndex(url, "/") + 1
	name := url[index:]
	index = strings.LastIndex(name, ".") + 1
	fileType := name[index:]
	return "E:\\file\\type_" + fileType + "\\" + name
}
