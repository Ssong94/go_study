package myapp

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/*
Go는 뒤에 _를 붙이고 test라고 하면 test코드로 작동한다.
ex) *_test.go
*/

/*
go get github.com/smartystreets/goconvey
github.com/stretchr/testify/assert
다운로드 후 테스트 하려는 디렉터리에 가서 Terminal -> goconvey 실행
*/

/*
테스트 코드는 함수명 앞에 Test라고 시작해야한다.
양식이 정해져있는데 t *testing.T -> testing 패키지에 T 포인터를 인자로 받아야한다.
*/
func TestIndexPathHandler(t *testing.T) {
	assert := assert.New(t)
	res := httptest.NewRecorder() // 실제 네트워크를 사용하지않고 실행가능
	req := httptest.NewRequest("GET", "/", nil)

	indexHandler(res, req) // app.go에서 핸들러를 분리해줌

	//if res.Code != http.StatusOK {
	//	t.Fatal("Failed!!", res.Code)
	//}

	assert.Equal(http.StatusOK, res.Code) // 잘 돌아가는지 아닌지 검사
}
func TestBarPathHandler_WithoutName(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello World!", string(data))
}

func TestBarPathHandler_WithName(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar?name=ssong", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello ssong!", string(data))
}

/* mux handler */
func TestIndexPathMuxHandler(t *testing.T) {
	assert := assert.New(t)
	res := httptest.NewRecorder() // 실제 네트워크를 사용하지않고 실행가능
	req := httptest.NewRequest("GET", "/", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello World", string(data))
}

const (
	jsonData = `{
"first_name":"ssong",
"last_name":"94",
"email":"ssong94@naver.com"
}`
)

// Json 테스트
func TestFooHandler_WithoutJson(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(jsonData))

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusCreated, res.Code)

	user := new(User)
	err := json.NewDecoder(res.Body).Decode(user)
	assert.Nil(err)
	assert.Equal("ssong", user.FirstName)
	assert.Equal("94", user.LastName)
	assert.Equal("ssong94@naver.com", user.Email)

}
