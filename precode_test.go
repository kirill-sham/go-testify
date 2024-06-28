package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе.
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	cafeExp := []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"}

	req := httptest.NewRequest("GET", "/cafe?count=5&city=moscow", nil)
	res := httptest.NewRecorder()

	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(res, req)

	nCafe := strings.Split(res.Body.String(), ",")
	lCafe := len(nCafe)

	assert.Equal(t, lCafe, totalCount)

	assert.Equal(t, nCafe, cafeExp)
}

func TestHTTPSuccessAndBodyNotAmty(t *testing.T) {
	// Запрос сформирован корректно, сервис возвращает код ответа 200
	baseUrl := url.Values{"count": []string{"3"}, "city": []string{"moscow"}}
	assert.HTTPSuccess(t, mainHandle, "GET", "/cafe", baseUrl)
	//тело не пустое
	assert.HTTPBody(mainHandle, "GET", "/cafe", baseUrl)
}

func TestMainHandlerWhenCiryNotFound(t *testing.T) {

	//Город, который передаётся в параметре city, не поддерживается. Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
	urlNotcity := url.Values{"count": []string{"10"}, "city": []string{"Notsity"}}
	assert.HTTPStatusCode(t, mainHandle, "GET", "/cafe", urlNotcity, http.StatusBadRequest)

	req := httptest.NewRequest("GET", "/cafe?count=4&city=Notsity", nil)
	res := httptest.NewRecorder()

	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(res, req)

	expected := `wrong city value`
	assert.Equal(t, res.Body.String(), expected)
}
