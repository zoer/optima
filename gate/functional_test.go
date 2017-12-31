package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGettingConfigParams(t *testing.T) {
	started := make(chan bool)
	done := make(chan bool)

	db, _ := dbConnect()
	defer db.Close()
	defer db.Exec("TRUNCATE configs CASCADE")

	pgConfigInsert(db, "test11", "param11", `{"x":"foo"}`)
	pgConfigInsert(db, "test11", "param22", `{"y":"foo"}`)

	var tsURL string

	go runService(func(handler http.Handler) error {
		ts := httptest.NewServer(handler)
		defer ts.Close()
		tsURL = ts.URL

		started <- true
		<-done
		return nil
	})

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("Unable to start the test server")
	}

	params := []struct {
		num    string
		body   io.Reader
		status int
		data   string
	}{
		{"1", newConfigRequestJSON("test11", "notexists"), 404, `{"error":"Configuration not found"}`},
		{"2", newConfigRequestJSON("test11", "param22"), 200, `{"y":"foo"}`},
		{"3", newConfigRequestJSON("test11", ""), 422, `{"error":"'Type' and 'Data' JSON values must be provided"}`},
		{"4", newConfigRequestJSON("", "param2"), 422, `{"error":"'Type' and 'Data' JSON values must be provided"}`},
		{"5", strings.NewReader(`<xml/>`), 400, `{"error":"Unable to read JSON body"}`},
	}
	for _, d := range params {
		t.Run(d.num, func(t *testing.T) {
			assert := assert.New(t)

			res, err := http.Post(tsURL+"/configs", "application/json", d.body)

			assert.NoError(err)
			assert.Equal(res.StatusCode, d.status)

			b, _ := ioutil.ReadAll(res.Body)
			res.Body.Close()
			assert.Equal(strings.TrimSpace(string(b)), d.data)
		})
	}

	done <- true
}
