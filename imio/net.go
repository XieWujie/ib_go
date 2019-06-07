package imio

import "net/http"

func RegisterHttpListener()  {
	http.Handle("/register",netHandler(handlerRegister))
	http.Handle("/login",netHandler(loginListener))
	http.Handle("/file/post",netHandler(handleFilePost))
	http.Handle("/addFriend",netHandler(handlerAddFriend))
	http.Handle("/file/get/",netHandler(handleFileGet))
	err := http.ListenAndServe(":8000",nil)
	handlerError(err)
}

