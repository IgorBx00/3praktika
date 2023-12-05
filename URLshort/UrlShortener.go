package main

import (
	"fmt"
	"html"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
)

func main() {
	fmt.Println("Сервер запущен. Ожидание подключений...")
	http.HandleFunc("/", handleHome)              // Обработчик запросов
	http.HandleFunc("/shorten", handleShortenURL) // Обработчик запросов
	http.HandleFunc("/red/", handleRedirect)      // Обработчик запросов
	http.ListenAndServe(":8080", nil)             // Запуск сервера на порту 8080
}

func getUrl(shortKey string) string {
	serv := ":6379"                    // берем адрес сервера из аргументов командной строки
	conn, err := net.Dial("tcp", serv) // открываем TCP-соединение к серверу
	if err != nil {
		fmt.Println(err) // вывод ошибки при открытии TCP-соединения к серверу
		return "Ошибка"
	}
	command := "HGET" + " " + shortKey + "//"
	if n, err := conn.Write([]byte(command)); n == 0 || err != nil {
		return "Erorr"
	}
	buff := make([]byte, 1024)
	var n int
	for i := 0; i < 3; i++ {
		n, err = conn.Read(buff) // получаем ответ
		if err != nil {
			return "Erorr"
		}
	}
	return string(buff[0 : n-1])
}

func Seturl(shortKey string, OrigUrl string) string {
	serv := ":6379"                    // берем адрес сервера из аргументов командной строки
	conn, err := net.Dial("tcp", serv) // открываем TCP-соединение к серверу
	if err != nil {
		fmt.Println(err) // вывод ошибки при открытии TCP-соединения к серверу
		return "Ошибка"
	}
	command := "HSET" + " " + shortKey + " " + OrigUrl + "//"
	if n, err := conn.Write([]byte(string(command))); n == 0 || err != nil {
		return "Erorr"
	}
	buff := make([]byte, 1024)
	var n int
	for i := 0; i < 4; i++ {
		n, err = conn.Read(buff) // получаем ответ
		if err != nil {
			break
		}
	}
	return string(buff[0 : n-1])
}

func CheckUrl(url string) string {
	serv := ":6379"                    // берем адрес сервера из аргументов командной строки
	conn, err := net.Dial("tcp", serv) // открываем TCP-соединение к серверу
	if err != nil {
		fmt.Println(err) // вывод ошибки при открытии TCP-соединения к серверу
		return "Ошибка"
	}
	command := "HURL" + " " + url + "//"
	if n, err := conn.Write([]byte(command)); n == 0 || err != nil {
		return "Erorr"
	}
	buff := make([]byte, 1024)
	var n int
	for i := 0; i < 3; i++ {
		n, err = conn.Read(buff) // получаем ответ
		if err != nil {
			break
		}
	}
	return string(buff[0 : n-1])
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost { // Если метод запроса - POST, перенаправляем на страницу /shorten
		http.Redirect(w, r, "/shorten", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html") // В противном случае, отображаем форму для ввода URL
	fmt.Fprint(w, `
	<!DOCTYPE html>
	<html>
	<head>
	<title>URL Shortener</title>
	</head>
	<body>
	<h2>URL Shortener</h2>
	<form method="post" action="/shorten">
	<input type="url" name="url" placeholder="Enter a URL" required>
	<input type="submit" value="Shorten">
	</form>
	</body>
	</html>
	`)
}

func generateShortURL(url string) string {
	bb := CheckUrl(url)
	if bb != "false" {
		return bb
	}
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortURL := ""
	for i := 0; i < 6; i++ {
		shortURL += string(letters[rand.Intn(len(letters))])
	}
	return shortURL
}

func handleShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Проверяем, что метод запроса - POST
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	originalURL := r.FormValue("url") // Получаем оригинальный URL из формы
	if originalURL == "" {            // Проверяем, что URL не пустой
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}
	shortKey := generateShortURL(html.UnescapeString(originalURL)) // Генерируем короткий ключ на основе оригинального URL
	fmt.Println(shortKey)
	Seturl(shortKey, html.UnescapeString(originalURL))
	shortenedURL := fmt.Sprintf("http://localhost:8080/red/%s", shortKey) // Формируем сокращенный URL
	log.Println(shortenedURL)
	// Serve the result page
	w.Header().Set("Content-Type", "text/html") // Отображаем страницу с результатами
	fmt.Fprint(w, `
	<!DOCTYPE html>
	<html>
	<head>
	<title>URL Shortener</title>
	</head>
	<body>
	<h2>URL Shortener</h2>
	<p>Original URL: `, originalURL, `</p>
	<p>Shortened URL: <a href="`, shortenedURL, `">`, shortenedURL, `</a></p>
	</body>
	</html>
	`)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := strings.TrimPrefix(r.URL.Path, "/red/") // Получаем сокращенный ключ из пути URL
	if shortKey == "" {                                 // Проверяем, что ключ не пустой
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}
	dbResponse := getUrl(shortKey) // Получаем оригинальный URL из базы данных
	if dbResponse == "not found" { // Проверяем, не найден ли ключ в базе данных
		http.Error(w, "Short key not found", 404)
		return
	}
	http.Redirect(w, r, string(dbResponse), http.StatusMovedPermanently) // Перенаправляем пользователя на оригинальный URL
}
