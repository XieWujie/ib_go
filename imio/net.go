package imio

import (
	"net/http"
)

func RegisterHttpListener() {
	http.Handle("/user/register", netHandler(handlerRegister))
	http.Handle("/user/login", netHandler(loginListener))
	http.Handle("/file/post", netWithToken(handleFilePost))
	http.Handle("/file/get/", netWithToken(handleFileGet))
	http.Handle("/user/find", netHandler(FindUser))
	http.Handle("/message/get", netHandler(requestMessageList))
	http.Handle("/relation/get", netHandler(requestUserRelation))
	http.Handle("/verify/get", netHandler(getVerify))
	http.Handle("/verify/send", netHandler(sendVerify))
	http.Handle("/room/create", netHandler(createRoom))
	http.Handle("/room/get", netHandler(getRoom))
	http.Handle("/user/getByIds", netHandler(findUserByIds))
	http.Handle("/user/getById", netHandler(findUserById))

	//http.HandleMsg("/pushService",websocket.Handler(LongConnect))
	StartWebsocket()
	err := http.ListenAndServe(":8000", nil)
	handlerError(err)
}
