package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (*http.Response, string) {
	t.Helper()

	var err error
	var req *http.Request
	var resp *http.Response
	var respBody []byte

	req, err = http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err = client.Do(req)
	require.NoError(t, err)

	respBody, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	require.NoError(t, err)

	return resp, string(bytes.TrimSpace(respBody))
}

func emptyMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func getTimeFromString(timeString string) time.Time {
	result, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		panic(err.Error())
	}
	return result
}

// func testGzippedRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (*http.Response, string) {
//	t.Helper()
//
//	var err error
//	var req *http.Request
//	var resp *http.Response
//	var respBody []byte
//
//	var b bytes.Buffer
//
//	w, _ := gzip.NewWriterLevel(&b, gzip.BestCompression)
//
//	if _, err = w.Write([]byte(body)); err != nil {
//		return nil, err.Error()
//	}
//
//	if err = w.Close(); err != nil {
//		return nil, err.Error()
//	}
//
//	req, err = http.NewRequest(method, ts.URL+path, bytes.NewReader(b.Bytes()))
//	require.NoError(t, err)
//
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Content-Encoding", "gzip")
//
//	client := &http.Client{
//		CheckRedirect: func(req *http.Request, via []*http.Request) error {
//			return http.ErrUseLastResponse
//		},
//	}
//
//	resp, err = client.Do(req)
//	require.NoError(t, err)
//
//	respBody, err = ioutil.ReadAll(resp.Body)
//	defer resp.Body.Close()
//
//	require.NoError(t, err)
//
//	return resp, string(bytes.TrimSpace(respBody))
//}
