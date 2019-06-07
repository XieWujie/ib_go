package imio

import (
	"../db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Friend struct {
	MarkName string `json:"mark_name"`
	FriendId string `json:"friend_id"`
	Group string `json:"group"`
}

func handlerAddFriend(w http.ResponseWriter, r * http.Request)*AppError{
	if r.Method == "POST"{
		if err := tokenVerify(r,w) ;err != nil{
			return err
		}
		var m map[string]string
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			return &AppError{Error:err,message:"解析requestBody出错",statusCode:500}
		}
		userid := m["userid"]
		friendid := m["friendid"]
		user,e := checkFriend(userid,friendid)
		if e != nil {
			return e
		}

		t,er := addFriend(&user.Friends,Friend{FriendId:friendid})
		if er != nil {
			return er
		}
		if err = user.Update("friends", t);err != nil{
			return &AppError{Error:err,message:"数据库存储失败",statusCode:500}
		}
		receipt := Receipt{StatusCode:http.StatusOK ,Description:"添加成功",Data:m}
		result,err := json.Marshal(receipt)
		_, _ = fmt.Fprintln(w, string(result))
	}else {
		return &AppError{message:"请求方式错误",statusCode:400}
	}
	return nil
}

func checkFriend(userid string,friendid string)(*db.User,*AppError)  {
	if len(userid)== 0 {
		return nil,&AppError{message:" userid不能为空",statusCode:500}
	}
	if len(friendid) == 0 {
		return nil,&AppError{message:" friendid不能为空",statusCode:500}
	}
	var user = db.User{UserId:friendid}
	err := user.Get()
	if err != nil {
		return nil,&AppError{Error:err,message:"请求数据库失败",statusCode:500}
	}
	if len(user.Username) == 0{
		return nil,&AppError{message:" friendid不存在",statusCode:500}
	}
	return &user,nil
}

func addFriend(elements *string,newElement Friend)(string,*AppError)  {
	if len(*elements)== 0 {
		new,err := json.Marshal(newElement)
		var s = "'["+string(new)+"]'"
		return s,&AppError{statusCode:500,message:"解析数据失败",Error:err}
	}
	var friends = make([]Friend,20)
	err := json.Unmarshal([]byte(*elements),&friends)
	if err != nil{
		log.Println(err)
		return "",&AppError{statusCode:500,message:"解析数据失败",Error:err}
	}
	for _,friend :=range friends{
		if friend.FriendId == newElement.FriendId {
			return "",&AppError{statusCode:400,message:"friendid已存在",Error:err}
		}
	}
	i := append(friends, newElement)
	target,er := json.Marshal(i)
	var s = "'"+string(target)+"'"
	if er != nil {
		return "",&AppError{statusCode:500,message:"解析数据失败",Error:err}
	}
	return s,nil
}