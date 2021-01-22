package service

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	tData  *Data
	tSaver *testSaver
	tCoder *testCoder
	tEcho  *echo.Echo
	tReq   *http.Request
	tRec   *httptest.ResponseRecorder
)

func initTest(t *testing.T) {
	tSaver = &testSaver{name: "test.wav"}
	tCoder = &testCoder{res: "olia"}
	tData = newTestData(tSaver, tCoder)
	tEcho = initRoutes(tData)
	tReq = httptest.NewRequest(http.MethodPost, "/convert", strings.NewReader(`{"audio":"aa"}`))
	tReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	tRec = httptest.NewRecorder()
}

func TestLive(t *testing.T) {
	initTest(t)
	req := httptest.NewRequest(http.MethodGet, "/live", nil)

	e := initRoutes(tData)
	e.ServeHTTP(tRec, req)
	assert.Equal(t, http.StatusOK, tRec.Code)
	assert.Equal(t, `{"service":"OK"}`, tRec.Body.String())
}

func TestValidate(t *testing.T) {
	d := &input{Data: "a"}
	assert.Nil(t, validate(d))
	d.Data = ""
	assert.NotNil(t, validate(d))
	d.Data = "d"
	d.Format = "mp3"
	assert.Nil(t, validate(d))
	d.Format = "m4a"
	assert.Nil(t, validate(d))
	d.Format = "maa"
	assert.NotNil(t, validate(d))
}

func TestConvert(t *testing.T) {
	initTest(t)

	tEcho.ServeHTTP(tRec, tReq)

	assert.Equal(t, http.StatusOK, tRec.Code)
	assert.Equal(t, `{"audio":"dGVzdA=="}`+"\n", tRec.Body.String())
}

func TestConvert_FailData(t *testing.T) {
	initTest(t)
	req := httptest.NewRequest(http.MethodPost, "/convert", strings.NewReader(`{"audio":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	
	tEcho.ServeHTTP(tRec, req)

	assert.Equal(t, http.StatusBadRequest, tRec.Code)
}

func TestConvert_FailType(t *testing.T) {
	initTest(t)
	req := httptest.NewRequest(http.MethodPost, "/convert", strings.NewReader(`{"audio":"a", "format":"maa"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	tEcho.ServeHTTP(tRec, req)

	assert.Equal(t, http.StatusBadRequest, tRec.Code)
}

func TestConvert_FailSaver(t *testing.T) {
	initTest(t)

	tSaver.err = errors.New("olia")

	tEcho.ServeHTTP(tRec, tReq)

	assert.Equal(t, http.StatusInternalServerError, tRec.Code)
}

func TestConvert_FailConvert(t *testing.T) {
	initTest(t)

	tCoder.err = errors.New("olia")

	tEcho.ServeHTTP(tRec, tReq)

	assert.Equal(t, http.StatusInternalServerError, tRec.Code)
}

func TestConvert_FailRead(t *testing.T) {
	initTest(t)

	tData.readFunc = func(string) ([]byte, error) { return nil, errors.New("olia") }

	tEcho.ServeHTTP(tRec, tReq)

	assert.Equal(t, http.StatusInternalServerError, tRec.Code)
}

type testSaver struct {
	name string
	err  error
	data bytes.Buffer
}

func (s *testSaver) Save(name string, reader io.Reader) (string, error) {
	io.Copy(&s.data, reader)
	return s.name, s.err
}

type testCoder struct {
	err  error
	name string
	res  string
}

func (s *testCoder) Convert(name string, format string, mt []string) (string, error) {
	s.name = name
	return s.res, s.err
}

func newTestData(s FileSaver, e Encoder) *Data {
	return &Data{Saver: s, Coder: e, readFunc: func(string) ([]byte, error) { return []byte("test"), nil }}
}
