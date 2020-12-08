package imio

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const OK = 200
const ERROR = 403

var Secret = []byte("imblog had login")

const LocalIp = "localhost:8000"

func init() {
	log.SetFlags(log.Llongfile)
}

type MgSend interface {
	Send(id string, message interface{}) error
}

type netHandler func(w http.ResponseWriter, r *http.Request) *AppError

type netWithToken func(w http.ResponseWriter, r *http.Request) *AppError

type postToken func(w http.ResponseWriter, r *http.Request) *AppError

type postHandler func(w http.ResponseWriter, r *http.Request) *AppError

type AppError struct {
	error      error
	message    string
	statusCode int
}

func (h netWithToken) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//err :=tokenVerify(r)
	//if err == nil{
	//	err = h(w,r)
	//}
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	err := h(w, r)
	if err != nil {
		message := err.message
		if len(message) == 0 {
			message = err.error.Error()
		}
		e := Receipt{StatusCode: err.statusCode, Description: message}
		u, _ := json.Marshal(e)
		http.Error(w, string(u), err.statusCode)
	}
}

func (h netHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
	r.ParseForm()
	if err := h(w, r); err != nil {
		message := err.message
		if len(message) == 0 {
			message = err.error.Error()
		}
		e := Receipt{StatusCode: err.statusCode, Description: message}
		u, _ := json.Marshal(e)
		http.Error(w, string(u), err.statusCode)
	}
}

func handlerError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Receipt struct {
	StatusCode  int         `json:"statusCode"`
	Data        interface{} `json:"data"`
	Description string      `json:"description"`
}

func createId() int {
	return autoId.Id()
}

type AutoInc struct {
	start, step int
	queue       chan int
	running     bool
}

var autoId *AutoInc = New(10000, 1)

func New(start, step int) (ai *AutoInc) {
	ai = &AutoInc{
		start:   start,
		step:    step,
		running: true,
		queue:   make(chan int, 4),
	}
	go ai.process()
	return

}

func (ai *AutoInc) process() {
	defer func() { recover() }()
	for i := ai.start; ai.running; i = i + ai.step {
		ai.queue <- i
	}
}

func (ai *AutoInc) Id() int {
	return <-ai.queue
}

func (ai *AutoInc) Close() {
	ai.running = false
	close(ai.queue)
}

func addElement(model []interface{}, elements *string, newElement interface{}) (string, error) {
	err := json.Unmarshal([]byte(*elements), model)
	if err != nil {
		return "", err
	}
	i := append(model, newElement)
	target, er := json.Marshal(i)
	return string(target), er
}

func errorReceipt(w http.ResponseWriter, code int, reason string) {
	receipt := Receipt{StatusCode: ERROR, Description: reason}
	result, _ := json.Marshal(receipt)
	_, _ = fmt.Fprintln(w, string(result))
}

func sendOk(w http.ResponseWriter) {
	receipt := &Receipt{StatusCode: 200, Description: "ok", Data: "ok"}
	rec, _ := json.Marshal(&receipt)
	fmt.Fprintln(w, string(rec))
}

func sendOkWithData(w http.ResponseWriter, data interface{}) {
	receipt := &Receipt{StatusCode: 200, Description: "ok", Data: data}
	rec, _ := json.Marshal(&receipt)
	fmt.Fprintln(w, string(rec))
	fmt.Println("push data", string(rec))
}
