package tests

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestALLPages(t *testing.T) {
	baseUrl := "http://localhost:3000"

	var tests = []struct {
		method   string
		url      string
		expected int
	}{
		{"GET", "/", 200},
		{"GET", "/about", 200},
		{"GET", "/notfound", 404},
		{"GET", "/articles", 200},
		{"GET", "/articles/create", 200},
		{"GET", "/articles/3", 200},
		{"GET", "/articles/3/edit", 200},
		{"POST", "/articles/3", 200},
		{"POST", "/articles", 200},
		{"POST", "/articles/1/delete", 404},
	}

	for _, test := range tests {
		t.Logf("当前请求 URL: %v \n", test.url)
		var (
			resp *http.Response
			err  error
		)
		if test.method == "POST" {
			data := make(map[string][]string)
			resp, err = http.PostForm(baseUrl+test.url, data)
		} else {
			resp, err = http.Get(baseUrl + test.url)
		}

		assert.NoError(t, err, "请求url:"+test.url+"报错")
		assert.Equal(t, test.expected, resp.StatusCode, "应返回状态码 "+strconv.Itoa(test.expected))
	}

}
