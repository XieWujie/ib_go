package imio

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)
import "../db"

func getVerify(w http.ResponseWriter, r *http.Request) *AppError {
	q := r.URL.Query()
	userId, _ := strconv.Atoi(q.Get("userId"))
	list := db.FindVerifyList(userId)
	sendOkWithData(w, list)
	return nil
}

func sendVerify(w http.ResponseWriter, r *http.Request) *AppError {
	var verify db.Verify
	_ = json.NewDecoder(r.Body).Decode(&verify)
	if verify.CreateAt == 0 {
		verify.CreateAt = time.Now().Unix()
	}
	if verify.State == db.NoAction {
		_ = verify.Save()
	} else if verify.State == db.Agree {
		conversationId := agreeAdd(verify)
		_ = verify.UpdateState()
		verify.VerifyInfo = strconv.FormatInt(int64(conversationId), 10)
	} else if verify.State == db.Defy {
		_ = verify.UpdateState()
	}

	sendOkWithData(w, verify)
	return nil
}
