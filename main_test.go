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

	for city, list := range cafeList {
		//Структуру объявил внутри цикла, чтобы сразу задать len(list) и не переопределять в процессе теста
		request := []struct {
			count int
			want  int
		}{
			{0, 0},
			{1, 1},
			{2, 2},
			{100, len(list)},
		}

		for _, v := range request {

			response := httptest.NewRecorder()
			req := httptest.NewRequest("GET",
				fmt.Sprintf("/cafe?city=%s&count=%d", city, v.count), nil)

			handler.ServeHTTP(response, req)

			require.Equal(t, http.StatusOK, response.Code,
				"city: %s, count: %d - status not Ok", city, v.count)

			cafes := strings.TrimSpace(response.Body.String())

			actualCount := 0

			if cafes != "" {
				actualCount = len(strings.Split(cafes, ","))
			}

			assert.Equal(t, v.want, actualCount,
				"city: %s, count: %d. result: get %d cafes:%v", actualCount, cafes)
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
		{"Завтрак", 1}, //для проверки регистра
	}

	for _, v := range request {

		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=moscow&search=%s", v.search), nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code,
			"Search: '%s' - invalid status", v.search)

		cafeStr := strings.TrimSpace(response.Body.String())

		var cafes []string

		if cafeStr != "" {
			cafes = strings.Split(cafeStr, ",")
		}
		//Проверяем количество кафе
		assert.Equal(t, v.wantCount, len(cafes),
			"Search: '%s' - expected %d, got %d", v.wantCount, len(cafes))

		searchLower := strings.ToLower(v.search)

		//Проверяем, что все кафе в списке - содержат поисковую строку
		for _, cafe := range cafes {

			assert.True(t, strings.Contains((strings.ToLower(cafe)), searchLower))

		}
	}

}
