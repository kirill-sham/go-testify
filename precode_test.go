package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	//Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе.
	moreCountUrl := url.Values{"count": []string{"10"}, "city": []string{"moscow"}}
	assert.HTTPBodyContains(t, mainHandle, "GET", "/cafe", moreCountUrl, "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент")
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
	assert.HTTPStatusCode(t, mainHandle, "GET", "/cafe", urlNotcity, 400)
	assert.HTTPBodyContains(t, mainHandle, "GET", "/cafe", urlNotcity, "wrong city value")
}
