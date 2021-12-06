package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadTest(t *testing.T) {
	assert := assert.New(t)

	path := "/home/kyugwang/Desktop/project/test.txt"
	file, _ := os.Open(path)
	defer file.Close() // defer는 함수가 종료되기 직전에 실행됨

	os.RemoveAll("./uploads") // 경로에 있는 파일 제거

	buf := &bytes.Buffer{}

	// NewWriter는 iowriter를 넘겨줘야하는데 buffer형식이다.
	writer := multipart.NewWriter(buf)                                      // 웹으로 파일 데이터를 전송할 때 사용하는 포멧, 데이터는 buf에 실려있다.
	multi, err := writer.CreateFormFile("upload_file", filepath.Base(path)) // filepath.Base => test.txt
	assert.NoError(err)                                                     // 이걸 통과하면 파일이 있다는 뜻
	io.Copy(multi, file)
	writer.Close()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/uploads", buf)
	req.Header.Set("Content-type", writer.FormDataContentType())

	uploadsHandler(res, req)
	assert.Equal(http.StatusOK, res.Code)

	/// 실제 폴더에 파일이 들어있는지 확인
	uploadFilePath := "./uploads/" + filepath.Base(path) // 업로드한 파일 경로
	_, err = os.Stat(uploadFilePath)                     // err가 위에서 이미 선언되었기 대문에 선언대입문(:=) 대신 일반대입문(=)을 사용한다.
	assert.NoError(err)                                  // 이걸 통과하면 파일이 있다는 뜻

	// 파일이 있다면 업로드된 파일하고 기존 파일하고 같은지 확인해야함.
	uploadFile, _ := os.Open(uploadFilePath) // 업로드한 파일 경로
	originFile, _ := os.Open(path)           // 기존 파일 경로
	defer uploadFile.Close()
	defer originFile.Close()

	uploadData := []byte{}
	originData := []byte{}
	uploadFile.Read(uploadData)
	originFile.Read(originData)

	assert.Equal(originData, uploadData) // 기존 파일과 업로드 파일이 같은지 검사

}
