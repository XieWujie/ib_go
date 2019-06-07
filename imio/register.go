package imio

import (
	"../db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)


func handlerRegister(w http.ResponseWriter, r * http.Request) *AppError {
	if r.Method == "POST"{
		var m map[string]string
		if err := json.NewDecoder(r.Body).Decode(&m);err != nil {
			return &AppError{Error:err,message:"处理requestBody失败",statusCode:500}
		}
		username := m["username"]
		password := m["password"]
		err := checkRegister(username,password)
		if err != nil{
			return err
		}
		id := createId(username,time.Now().String())
		m["token"],err = createToken(id)
		if err != nil{
			return err
		}
		receipt := Receipt{StatusCode:OK,Description:"注册成功",Data:m}
		result,e := json.Marshal(receipt)
		if e != nil {
			return &AppError{Error:e,message:e.Error(),statusCode:500}
		}
		_, _ = fmt.Fprintln(w, string(result))
		user := db.User{UserId:id,Username:username,Password:password}
		e = user.Save()
		if e != nil {
			return &AppError{Error:e,message:"保存用户信息失败,请重新注册",statusCode:500}
		}
	}else {
		errorReceipt(w,ERROR,"请求方式错误")
	}
	return nil
}



func checkRegister(username string,password string)*AppError{
	if len(username)== 0 {
		return  &AppError{message:"用户名不能为空",statusCode:http.StatusBadRequest}
	}
	if len(password) == 0 {
		return  &AppError{message:"密码不能为空",statusCode:http.StatusBadRequest}
	}
	var user = db.User{Username:username}
	err := user.Get()
	if err != nil {
		log.Fatal(err)
	}
	if len(user.UserId) >0  {
		return &AppError{message:"账号已经被注册",statusCode:400}
	}
	return  nil
}

