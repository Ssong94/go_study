package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/uploads", uploadsHandler)
	http.Handle("/", http.FileServer(http.Dir("public")))

	http.ListenAndServe(":3000", nil)
}

func uploadsHandler(w http.ResponseWriter, r *http.Request) {
	uploadFile, header, err := r.FormFile("upload_file") // id가 upload_file

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	// defer는 함수가 종료되기 직전에 실행됨
	defer uploadFile.Close() // 파일을 만들고 닫아줘야함(os자원이라 반납해야함)

	dirName := "./uploads"
	os.MkdirAll(dirName, 0777)                                 // dirname 폴더가 없으면 만들어줌, 777 -> read,write,execute 가능
	filePath := fmt.Sprintf("%s/%s", dirName, header.Filename) // 폴더명/파일명, 파일명은 header에 들어있다.
	file, err := os.Create(filePath)                           // 비어있는 새로운 파일을 만듬
	defer file.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	io.Copy(file, uploadFile)    // 비어있는 파일에 uploadFile을 복사해준다.
	w.WriteHeader(http.StatusOK) // err없이 잘 되면 OK신호
	fmt.Fprint(w, filePath)      // 어디에 업로드되는지 출력

}
