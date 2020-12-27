package imio

import (
	"../db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func changeCustomFriendBg(w http.ResponseWriter, r *http.Request) *AppError {
	_ = r.ParseMultipartForm(32 << 20)
	var message = r.FormValue("json")
	var m map[string]int
	_ = json.Unmarshal([]byte(message), &m)
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		return &AppError{statusCode: 400, error: err}
	}
	defer file.Close()
	var filePath = createFilePath(handler.Filename)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		log.Println(err)
		return &AppError{statusCode: 400, error: err}
	}
	url := createUrl(handler.Filename)
	var user = db.User{UserId:m["ownerId"]}
	_ = user.Get()
	conversationId := m["conversationId"]
	for i, v := range user.Friends {
		if v.ConversationId == conversationId {
			user.Friends[i].Background = url
			break
		}
	}
	_ = user.Update()
	var newM = make(map[string]string)
	newM["background"] = url
	sendOkWithData(w, newM)
	return nil
}

type markNameEntity struct {
	OwnerId        int `json:"ownerId"`
	ConversationId int `json:"conversationId"`
	MarkName     string `json:"markName"`
}

func upDateFriendMarkName(w http.ResponseWriter,r * http.Request)*AppError  {
	var en  = new(markNameEntity)
	_ = json.NewDecoder(r.Body).Decode(&en)
	var owner = db.User{UserId:en.OwnerId}
	_ = owner.Get()
	var isUpdate = false
	for i,v := range owner.Friends{
		if en.ConversationId == v.ConversationId{
			owner.Friends[i].MarkName = en.MarkName
			_ = owner.Update()
			sendOk(w)
			isUpdate = true
			break
		}
	}
	if !isUpdate{
		return &AppError{statusCode:400,message:"找不到conversationId"}
	}
	return nil
}

func agreeAdd(m db.Verify) int {

	member := make([]db.MemberInfo, 2)
	members := append(member, db.MemberInfo{UserId: m.UserFrom}, db.MemberInfo{UserId: m.UserTo})
	cov := db.Conversation{ConversationType: db.COV_TYPE_SINGLE, Members: members}
	_ = cov.Save()
	user := &db.User{UserId: m.UserFrom}
	_ = user.Get()
	user.Friends = append(user.Friends, db.RelationShip{UserId: m.UserTo, ConversationId: cov.ConversationId,Notify:true})
	_ = user.Update()
	friend := &db.User{UserId: m.UserTo}
	_ = friend.Get()
	friend.Friends = append(friend.Friends, db.RelationShip{UserId: m.UserFrom, ConversationId: cov.ConversationId,Notify:true})
	_ = friend.Update()
	return cov.ConversationId
}


func friendMsgNotify(w http.ResponseWriter, r *http.Request) *AppError {
	en := new(notifyEntity)
	_ = json.NewDecoder(r.Body).Decode(&en)
	user := db.User{UserId: en.OwnerId}
	_ = user.Get()
	for i, v := range user.Friends {
		if v.ConversationId == en.ConversationId {
			user.Friends[i].Notify = en.Notify
			break
		}
	}
	_ = user.Update()
	sendOk(w)
	return nil
}

func requestUserRelation(w http.ResponseWriter, r *http.Request) *AppError {
	fmt.Println("request user relationship")
	q := r.URL.Query()
	ownerId, _ := strconv.Atoi(q.Get("userId"))
	owner := db.User{UserId: ownerId}
	_ = owner.Get()
	var relations = owner.Friends
	list := make([]map[string]interface{}, len(relations))
	for i, v := range relations {
		m := make(map[string]interface{})
		user := db.User{UserId: v.UserId}
		user.Get()
		m["user"] = user
		m["conversationId"] = v.ConversationId
		m["notify"] = v.Notify
		m["markName"] = v.MarkName
		m["background"] = v.Background
		list[i] = m
	}
	sendOkWithData(w, list)
	return nil
}