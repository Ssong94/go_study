package myapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	FirstName string    `json:"first_name"` // json alias
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type fooHandler struct {
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello World")
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!", name)
}

func NewHttpHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	})
	mux.HandleFunc("/bar", barHandler)
	mux.Handle("/foo", &fooHandler{})
	//http.ListenAndServe(":3000", mux)
	return mux
}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user) // NewDecoder는 io.Reader를 인자로 받고 Body는 io.Reader를 포함하고 있다.
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // hear로 오류가 있다는 것을 알려줌
		fmt.Fprint(w, "Bad Request: ", err)  // body에 에러를 알려줌
		return
	}
	user.CreatedAt = time.Now()

	data, _ := json.Marshal(user) // 인터페이스를 받아서 json형태로 바꿔주는 메소드(byte와 err를 리턴함)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(data))
}
