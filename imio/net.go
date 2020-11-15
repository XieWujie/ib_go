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
	http.Handle("/user/relation", netHandler(relationWithUser))
	http.Handle("/relation/get", netHandler(requestUserRelation))
	//http.HandleMsg("/pushService",websocket.Handler(LongConnect))
	StartWebsocket()
	err := http.ListenAndServe(":8000", nil)
	handlerError(err)
}
