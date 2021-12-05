package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type user struct {
	FirstName string `json:"first_name"` // json alias
	LastName  string `json:"last_name"`
	Email     string
	CreatedAt time.Time
}

type fooHandler struct {
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!", name)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	})

	mux.HandleFunc("/bar", barHandler)

	mux.Handle("/foo", &fooHandler{})

	http.ListenAndServe(":3000", mux)
}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := new(user)
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
