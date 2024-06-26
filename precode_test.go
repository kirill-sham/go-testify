package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {

	// Запрос сформирован корректно, сервис возвращает код ответа 200
	baseUrl := url.Values{"count": []string{"3"}, "city": []string{"moscow"}}
	assert.HTTPSuccess(t, mainHandle, "GET", "/cafe", baseUrl)

	//тело не пустое
	assert.HTTPBody(mainHandle, "GET", "/cafe", baseUrl)

	//Город, который передаётся в параметре city, не поддерживается. Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
	urlNotcity := url.Values{"count": []string{"10"}, "city": []string{"Notsity"}}
	assert.HTTPStatusCode(t, mainHandle, "GET", "/cafe", urlNotcity, 400)
	assert.HTTPBodyContains(t, mainHandle, "GET", "/cafe", urlNotcity, "wrong city value")

	//Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе.
	moreCountUrl := url.Values{"count": []string{"10"}, "city": []string{"moscow"}}
	assert.HTTPBodyContains(t, mainHandle, "GET", "/cafe", moreCountUrl, "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент")

}
