package myapp

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello World", string(data))

}

func TestUsers(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	res, err := http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	data, _ := ioutil.ReadAll(res.Body)
	assert.Contains(string(data), "Get UserInfo")
}

func TestGetUserInfo(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	res, err := http.Get(ts.URL + "/users/89")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	data, _ := ioutil.ReadAll(res.Body)
	assert.Contains(string(data), "No User Id:89")

	res2, err := http.Get(ts.URL + "/users/51")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res2.StatusCode)

	data2, _ := ioutil.ReadAll(res2.Body)
	assert.Contains(string(data2), "No User Id:51")
}

const (
	jsonData = `{
"first_name":"ssong",
"last_name":"94",
"email":"ssong94@naver.com"
}`
)

func TestCreateUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/users", "application/json", strings.NewReader(jsonData))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	id := user.ID
	resp, err = http.Get(ts.URL + "/users/" + strconv.Itoa(id))
	assert.NoError(err)
	assert.Equal(resp.StatusCode, http.StatusOK)

	user2 := new(User)
	err = json.NewDecoder(resp.Body).Decode(user2)
	assert.NoError(err)
	assert.Equal(user.ID, user2.ID) // 서로 같아야 맞음
	assert.Equal(user.FirstName, user2.FirstName)

}
