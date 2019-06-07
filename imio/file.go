package imio

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func handleFilePost(w http.ResponseWriter, r *http.Request) *AppError {
	if r.Method == "POST" {
		if err := tokenVerify(r,w) ;err != nil{
			return err
		}
		r.ParseForm()
		fileType := r.Form.Get("fileType")
		if len(fileType) == 0 {
			return &AppError{statusCode: 400, message: "fileType不能为空"}
		}
		filePath,fileName := GetFile(fileType)
		f, e := os.Create(filePath)
		if e != nil {
			return &AppError{Error: e, statusCode: 500}
		}
		_, err := io.Copy(f, r.Body)
		if err != nil {
			return &AppError{Error: err, statusCode: 500}
		}
		m := make(map[string]string)
		m["url"] = createUrl(fileName)
		rp := Receipt{StatusCode: http.StatusOK, Description: "上传成功", Data: m}
		t, e := json.Marshal(rp)
		if e != nil {
			return &AppError{Error: e, statusCode: 500}
		}
		_, err = fmt.Fprintln(w, string(t))
		defer f.Close()
	} else {
		return &AppError{message: "请求方式错误", statusCode: 400}
	}
	return nil
}

func handleFileGet(w http.ResponseWriter, r *http.Request) *AppError  {
	if r.Method == "GET" {
		filePath := getPathFromUrl(r.RequestURI)
		println(filePath)
		file,err := os.Open(filePath)
		if err != nil {
			return &AppError{statusCode:400,message:"找不到文件",Error:err}
		}
		_, err = io.Copy(w, file)
		defer file.Close()
	}
	return nil
}

func GetFile(fileType string) (string,string){
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	len := 10
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	fileName := string(bytes)+ "." + fileType
	basePath := "F:\\file\\type_" + fileType
	_ = os.MkdirAll(basePath,os.ModeAppend)
	filePath := basePath +"\\"+fileName
	return filePath,fileName
}
func createUrl(fileName string)string  {
	url:= LOCAL_IP+"/file/get/"+fileName
	return url
}

func getPathFromUrl(url string)string  {
	index := strings.LastIndex(url,"/")+1
	name := url[index:]
	index = strings.LastIndex(name,".")+1
	fileType := name[index:]
	return "F:\\file\\type_" + fileType+"\\"+name
}
