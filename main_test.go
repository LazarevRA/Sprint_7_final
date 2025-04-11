package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)

	}

}

func TestCafeWhenNotOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}

	for _, v := range requests {

		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}

}

func TestCafeCount(t *testing.T) {

	handler := http.HandlerFunc(mainHandle)

	request := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, 0},
	}

	for city, _ := range cafeList {
		for _, v := range request {

			response := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=%s&count=%d", city, v.count), nil)

			handler.ServeHTTP(response, req)

			require.Equal(t, http.StatusOK, response.Code)

			cafes := strings.Split(response.Body.String(), ",")

			actualCount := len(cafes)

			if cafes[0] == "" {
				actualCount = 0
			}

			//if actualCount == 1 {
			//	cafes[0] = strings.TrimSpace(cafes[0]) // we can use trimspace, if StringSplit for empty string != ""
			//	if cafes[0] == "" {
			//		actualCount = 0
			//	}
			//}

			if v.count > actualCount {
				v.want = len(cafeList[city])
			}
			assert.Equal(t, v.want, actualCount)
		}
	}
}

func TestCafeSearch(t *testing.T) {

	handler := http.HandlerFunc(mainHandle)

	request := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range request {

		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=moscow&search=%s", v.search), nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		cafes := strings.Split(response.Body.String(), ",")

		//actualCount := len(cafes)

		//if cafes[0] == "" {
		//	actualCount = 0
		//}
		actualCount := 0

		for _, cafe := range cafes {
			if strings.Contains(strings.ToLower(cafe), v.search) {
				actualCount++
			}
		}
		assert.Equal(t, v.wantCount, actualCount)
	}

}
