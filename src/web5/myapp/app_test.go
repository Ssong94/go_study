package myapp

import (
	"encoding/json"
	"fmt"
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
	//assert.Contains(string(data), "Get UserInfo")
	assert.Contains(string(data), "No Users")
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
	assert.Equal(user.ID, user2.ID) // ?????? ????????? ??????
	assert.Equal(user.FirstName, user2.FirstName)

}

func TestDeleteUser(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// ??????
	resp, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(jsonData))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	// ??????
	req, _ := http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(resp.StatusCode, http.StatusOK)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "Deleted User ID:1")
}

const (
	updateJsonData = `{"id: 1"}`
)

func TestUpdateUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(jsonData))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	// create??? ID??? ????????? ??????
	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	updateStr := fmt.Sprintf(`{"id":%d, "first_name":"updated"}`, user.ID) // ???????????? ID??? ?????????
	// updateStr??? ???????????? ?????? user??? ????????? ???????????????.
	req, _ := http.NewRequest("PUT", ts.URL+"/users",
		strings.NewReader(updateStr))
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// update??? ?????? ????????? ????????? ????????? ????????? ??????
	updateUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(updateUser)
	assert.NoError(err)

	assert.Equal(updateUser.ID, user.ID)             // create??? ID??? update??? ID??? ????????? ??????.
	assert.Equal("updated", updateUser.FirstName)    // create??? ??? update??? FirstName??? update??? ???????????? ????????? ??????.
	assert.Equal(user.LastName, updateUser.LastName) // update ????????? LastName??? ????????? ??????.
	assert.Equal(user.Email, updateUser.Email)       // update ????????? Email??? ????????? ??????.

}

func TestUsers_WithUsersData(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// ?????? 2???
	res, err := http.Post(ts.URL+"/users", "application/json", strings.NewReader(jsonData))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)
	res, err = http.Post(ts.URL+"/users", "application/json", strings.NewReader(jsonData))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	// ?????? ??????
	res, err = http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	var users []*User
	err = json.NewDecoder(res.Body).Decode(&users)
	assert.NoError(err)
	assert.Equal(2, len(users))

}
