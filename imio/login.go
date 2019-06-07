package imio

import (
	"../db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func loginListener(w http.ResponseWriter, r * http.Request) *AppError{
	 if r.Method == "POST"{
		var m map[string]string
		 if err := json.NewDecoder(r.Body).Decode(&m);err != nil {
			 return &AppError{Error:err,message:"处理requestBody失败",statusCode:500}
		 }
		username := m["username"]
		password := m["password"]
		user,err := checkLogin(username,password)
		 if err != nil{
			 return err
		 }
		m["token"],err = createToken(user.UserId)
		 if err != nil{
			 return err
		 }
		receipt := Receipt{StatusCode:OK,Description:"登陆成功",Data:m}
		result,e := json.Marshal(receipt)
		 if e != nil {
			 return &AppError{Error:e,message:e.Error(),statusCode:500}
		 }
		_, _ = fmt.Fprintln(w, string(result))
	}else {
		return &AppError{message:"请求方式错误",statusCode:400}
	 }
	 return nil
}

func checkLogin(username string,password string)( *db.User,*AppError ) {
	if len(username)== 0 {
		return nil, &AppError{message:"用户名不能为空",statusCode:http.StatusBadRequest}
	}
	if len(password) == 0 {
		return nil, &AppError{message:"密码不能为空",statusCode:http.StatusBadRequest}
	}
	var user = db.User{Username:username}
	err := user.Get()
	if err != nil {
		log.Fatal(err)
	}
	if user.Password != password{
		return nil, &AppError{message:"账号或密码错误",statusCode:http.StatusBadRequest}
	}
	return  &user,nil
}