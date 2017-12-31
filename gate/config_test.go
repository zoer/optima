package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigService_HandleParams(t *testing.T) {
	repo := NewInMemoryConfigRepo()
	repo.SetParams("test1", "param1", []byte(`{"x": "foo"}`))
	repo.SetParams("test1", "param2", []byte(`{"y": "boo"}`))
	repo.SetParams("test2", "param1", []byte(`{"z": "baz"}`))

	service := NewConfigService(repo)

	params := []struct {
		num, configKey, paramKey string
		status                   int
		body                     string
	}{
		{"1", "test1", "param1", 200, `{"x": "foo"}`},
		{"2", "test1", "param2", 200, `{"y": "boo"}`},
		{"3", "test2", "param1", 200, `{"z": "baz"}`},
		{"4", "test1", "notexists", 404, `{"error":"Configuration not found"}`},
	}

	for _, d := range params {
		t.Run(d.num, func(t *testing.T) {
			assert := assert.New(t)

			r := newConfigRequest(d.configKey, d.paramKey)
			w := httptest.NewRecorder()
			service.HandleParams(w, r)

			assert.Equal(w.Code, d.status)
			assert.Equal(strings.TrimSpace(string(w.Body.Bytes())), d.body)
		})
	}
}

func newConfigRequest(config, param string) *http.Request {
	r, _ := http.NewRequest("POST", "/configs", newConfigRequestJSON(config, param))
	return r
}

func newConfigRequestJSON(config, param string) io.Reader {
	j := fmt.Sprintf(`{"Type": "%s", "Data": "%s"}`, config, param)
	return strings.NewReader(j)
}
