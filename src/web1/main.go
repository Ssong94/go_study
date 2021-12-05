package main

import (
	"fmt"
	"net/http"
)

type fooHandler struct {
}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Bar!")
}

func barHandler2(w http.ResponseWriter, r *http.Request) {
	// url에서 argument 뽑아냄
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "wow"
	}
	fmt.Fprintf(w, "정말 신기하군요 %s !", name)
}

func main() {
	/* 함수를 직접 등록 */
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	})
	http.HandleFunc("/bar", barHandler)

	/* 인스턴스 형태로 등록(인스턴스를 만들고 거기에 해당하는 인터페이스를 구현) */
	http.Handle("/foo", &fooHandler{})

	/* 위에 테스트 하려면 밑에 주석 해제 */
	//http.ListenAndServe(":3000", nil)

	/* 새로운 라우터 인스턴스를 만들어서 그 인스턴스를 넘겨주는 방식으로 구현 */
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world")
	})

	mux.HandleFunc("/bar", barHandler2)

	http.ListenAndServe(":3030", mux)

}
