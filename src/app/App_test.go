package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"shorterUrl/src/app"
	"shorterUrl/src/env"
	"testing"
)

const (
	expTime       = 60
	longUrl       = "https://www.baidu.com"
	shortLink     = "IFHzaO"
	shortLinkInfo = `{
"url":"https://www.baidu.com.com",
"created_at":"2020-03-29 15:27:07.202571 +0800 CST m=+48.478071708",
"expiration_in_minutes":60}`
)

type storageMock struct {
	mock.Mock
}

var testApp app.App
var mockR *storageMock

func (s *storageMock) Shorten(url string, exp int64) (string, error) {
	args := s.Called(url, exp)
	return args.String(0), args.Error(1)
}

func (s *storageMock) ShortLinkInfo(eid string) (interface{}, error) {
	args := s.Called(eid)
	return args.String(0), args.Error(1)
}

func (s *storageMock) UnShorten(eid string) (string, error) {
	args := s.Called(eid)
	return args.String(0), args.Error(1)
}

func init() {
	testApp = app.App{}
	mockR = new(storageMock)
	testApp.Initialize(&env.Env{S: mockR})
}

func TestCreateShortLink(t *testing.T) {
	jsonStr := []byte(`{
		"url":"https://www.baidu.com",
		"expiration_in_minutes":60
	}`)

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatalf("can't create a request %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	mockR.On("Shorten", longUrl, int64(expTime)).Return(shortLink, nil).Once()
	rw := httptest.NewRecorder()
	testApp.Router.ServeHTTP(rw, req)
	if rw.Code != http.StatusCreated {
		t.Fatalf("Expected %d ,got %d", http.StatusCreated, rw.Code)
	}
	resp := struct {
		ShortLink string `json:"shortlink"`
	}{}

	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatal("can't decode response")
	}
	if resp.ShortLink != shortLink {
		t.Fatalf("Expected %s, got %s.", shortLink, resp.ShortLink)
	}
}

func TestRedirect(t *testing.T) {
	r := fmt.Sprintf("/%s", shortLink)
	req, err := http.NewRequest("GET", r, nil)
	if err != nil {
		t.Fatalf("can't create a request %v", err)
	}
	mockR.On("UnShorten", shortLink).Return(longUrl, nil).Once()
	rw := httptest.NewRecorder()
	testApp.Router.ServeHTTP(rw, req)
	if rw.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected %d , got %d .", http.StatusTemporaryRedirect, rw.Code)
	}
}
