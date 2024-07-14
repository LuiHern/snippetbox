package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.betocodes.io/internal/assert"
)

func TestPing(t *testing.T) {
	// init response recorder
	rr := httptest.NewRecorder()

	// init dummy http req
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// call ping handler func, pass the recorder and request
	ping(rr, r)

	// call result to get the response from the ping
	rs := rr.Result()

	// assert statusCode
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// check the body
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
