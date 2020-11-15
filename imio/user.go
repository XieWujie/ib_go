package imio

import (
	"../db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func loginListener(w http.ResponseWriter, r *http.Request) *AppError {
	if r.Method == "POST" {
		var m map[string]string
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			return &AppError{error: err, message: "处理requestBody失败", statusCode: 500}
		}
		username := m["username"]
		password := m["password"]
		user, err := checkLogin(username, password)
		if err != nil {
			return err
		}
		rec := make(map[string]interface{})
		rec["token"], _ = createToken(string(user.UserId))
		rec["userId"] = user.UserId
		rec["username"] = username
		receipt := Receipt{StatusCode: OK, Description: "ok", Data: rec}
		result, e := json.Marshal(receipt)
		if e != nil {
			return &AppError{error: e, message: e.Error(), statusCode: 500}
		}
		_, _ = fmt.Fprintln(w, string(result))
	} else {
		return &AppError{message: "请求方式错误", statusCode: 400}
	}
	return nil
}

func checkLogin(username string, password string) (*db.User, *AppError) {
	if len(username) == 0 {
		return nil, &AppError{message: "用户名不能为空", statusCode: http.StatusBadRequest}
	}
	if len(password) == 0 {
		return nil, &AppError{message: "密码不能为空", statusCode: http.StatusBadRequest}
	}
	var user = db.User{Username: username}
	err := user.Get()
	if err != nil {
		log.Fatal(err)
	}
	if user.Password != password {
		return nil, &AppError{message: "账号或密码错误", statusCode: http.StatusBadRequest}
	}
	return &user, nil
}

func handlerRegister(w http.ResponseWriter, r *http.Request) *AppError {
	if r.Method == "POST" {
		var m map[string]string
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			return &AppError{error: err, message: "处理requestBody失败", statusCode: 500}
		}
		username := m["username"]
		password := m["password"]
		err := checkRegister(username, password)
		if err != nil {
			return err
		}

		user := db.User{Username: username, Password: password}
		e := user.Save()
		if e != nil {
			fmt.Println(e)
			return &AppError{error: e, message: "保存用户信息失败,请重新注册", statusCode: 500}
		}
		rec := make(map[string]string)
		rec["token"], _ = createToken(username)
		rec["userid"] = strconv.Itoa(user.UserId)
		rec["username"] = username
		receipt := Receipt{StatusCode: OK, Description: "注册成功", Data: rec}
		result, _ := json.Marshal(receipt)
		fmt.Println(string(result))
		_, _ = fmt.Fprintln(w, string(result))
	} else {
		errorReceipt(w, ERROR, "请求方式错误")
	}
	return nil
}

func checkRegister(username string, password string) *AppError {
	if len(username) == 0 {
		return &AppError{message: "用户名不能为空", statusCode: http.StatusBadRequest}
	}
	if len(password) == 0 {
		return &AppError{message: "密码不能为空", statusCode: http.StatusBadRequest}
	}
	var user = db.User{Username: username}
	err := user.Get()
	if err != nil {
		log.Fatal(err)
	}
	if user.UserId > 0 {
		return &AppError{message: "账号已经被注册", statusCode: 400}
	}
	return nil
}

func FindUser(w http.ResponseWriter, r *http.Request) *AppError {
	q := r.URL.Query()
	key := q.Get("key")
	if len(key) == 0 {
		return &AppError{statusCode: 403, message: "key 为空"}
	}
	list := db.UserSearch(key)
	receipt := Receipt{StatusCode: OK, Description: "ok", Data: list}
	rec, _ := json.Marshal(receipt)
	body := string(rec)
	fmt.Fprint(w, body)
	return nil
}

func relationWithUser(w http.ResponseWriter, r *http.Request) *AppError {
	q := r.URL.Query()
	ownerId, _ := strconv.Atoi(q.Get("ownerId"))
	userId, _ := strconv.Atoi(q.Get("userId"))
	owner := db.User{UserId: ownerId}
	owner.Get()
	var relations []int
	_ = json.Unmarshal([]byte(owner.Friends), &relations)
	isFriend := false
	for _, v := range relations {
		if v == userId {
			isFriend = true
			break
		}
	}
	m := make(map[string]interface{})
	m["isFriend"] = isFriend
	sendOkWithData(w, m)
	return nil
}

func requestUserRelation(w http.ResponseWriter, r *http.Request) *AppError {
	fmt.Println("request user relationship")
	q := r.URL.Query()
	ownerId, _ := strconv.Atoi(q.Get("userId"))
	owner := db.User{UserId: ownerId}
	owner.Get()
	var relations []db.RelationShip
	_ = json.Unmarshal([]byte(owner.Friends), &relations)
	list := make([]map[string]interface{}, len(relations))
	for i, v := range relations {
		m := make(map[string]interface{})
		user := db.User{UserId: v.UserId}
		user.Get()
		m["user"] = user
		m["conversationId"] = v.ConversationId
		list[i] = m
	}
	sendOkWithData(w, list)
	return nil
}
