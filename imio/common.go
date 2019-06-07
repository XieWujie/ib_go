package imio

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"log"
	"net/http"
	"time"
)

const OK = 200
const ERROR = 403
var Secret  = []byte("imblog had login" )

const LOCAL_IP  = "localhost:8000"
func init() {
	log.SetFlags(log.Llongfile)
}
type MgSend interface {
	Send(id string, message interface{}) error
}

type netHandler func(w http.ResponseWriter, r * http.Request) *AppError

type AppError struct {
	Error error
	message string
	statusCode int
}

func (h netHandler)ServeHTTP(w http.ResponseWriter, r * http.Request)  {
	if err := h(w,r);err != nil{
		message := err.message
		if len(message) == 0{
			message = err.Error.Error()
		}
		e := Receipt{StatusCode:err.statusCode,Description:message}
		u,_ := json.Marshal(e)
		http.Error(w,string(u),err.statusCode)
	}
}

type TokenError struct {
	message string
}

func (err *TokenError) Error() string {
	return err.message
}

func handlerError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Receipt struct {
	StatusCode  int       `json:"statusCode"`
	Data        interface{} `json:"data"`
	Description string      `json:"description"`
}
type IBClaims struct {
	id string
	jwt.StandardClaims
}

func createToken(key string)(string, *AppError) {
	claims := IBClaims{
		key,
		jwt.StandardClaims{
			ExpiresAt:time.Now().Add(time.Hour*time.Duration(72)).Unix(),
			IssuedAt:time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	signedToken,err := token.SignedString(Secret)
	if err != nil {
		log.Fatal(err)
		return "",&AppError{Error:err,message:"生成token失败",statusCode:500}
	}
	return signedToken,nil
}

func createId(keys...string)string{
	md := md5.New()
	var b bytes.Buffer
	for _,key := range keys  {
		b.WriteString(key)
	}
	md.Write(b.Bytes())
	return string(md.Sum(nil))
}

func addElement(model []interface{},elements *string,newElement interface{})(string,error)  {
	err := json.Unmarshal([]byte(*elements),model)
	if err != nil{
		return "",err
	}
	i := append(model, newElement)
	target,er := json.Marshal(i)
	return string(target),er
}

func tokenVerify(r *http.Request,w http.ResponseWriter)*AppError{
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(Secret), nil
		})
	if err == nil {
		if !token.Valid {
			return &AppError{message:"Token is not valid",statusCode:400}
		}
	}else{
		return &AppError{Error:err,message:err.Error(),statusCode:500}
	}
	return nil
}

func errorReceipt(w http.ResponseWriter, code int, reason string) {
	receipt := Receipt{StatusCode: ERROR, Description: reason}
	result, _ := json.Marshal(receipt)
	_, _ = fmt.Fprintln(w, string(result))
}

