// +build integration

// Adapted from: https://github.com/shurcooL/httpgzip. Thanks sir!

package jgoweb

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Test that GzipServeContent correctly determines the content type as "text/plain",
// not as "application/x-gzip".
func TestGzipServeContentDetectContentType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		content := "This is some plain text that compresses easily. " +
			strings.Repeat("NaN", 16) + " Batman!"

		GzipServeContent(w, req, "", time.Time{}, strings.NewReader(content))
	}))
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	got := resp.Header.Get("Content-Type")
	want := "text/plain; charset=utf-8"
	if got != want {
		t.Errorf("got:\n%v\nwant:\n%v\n", got, want)
	}
}

// Test that if the handler already explicitly set "Content-Encoding" header,
// then ServeContent shouldn't try to do apply compression, just serve as is.
func TestGzipServeContentExplicitContentEncoding(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		content := "This is some plain text that compresses easily. " +
			strings.Repeat("NaN", 16) + " Batman!"

		w.Header()["Content-Encoding"] = nil
		GzipServeContent(w, req, "", time.Time{}, strings.NewReader(content))
	}))
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	got := resp.Header.Get("Content-Encoding")
	want := ""
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q\n", got, want)
	}
}
