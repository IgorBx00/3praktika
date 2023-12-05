package main

import (
	"fmt"
	"net"
	"time"
)

// СТЕК
var stack = &Stack{}

// ОЧЕРЕДЬ
var queue = &Queue{}

// ХЕШ-ТАБЛИЦА
var hashtable = &HashMap{}

// МНОЖЕСТВО
var set = &SetMap{}

func main() {
	listener, _ := net.Listen("tcp", ":6379") // открываем слушающий сокет
	fmt.Println("Сервер запущен. Ожидание подключений...")
	for {
		conn, err := listener.Accept() // принимаем TCP-соединение от клиента и создаем новый сокет
		if err != nil {
			continue // при возникновении ошибки у 1 из соединений программа не закрывается и не обрабатывает запрос от этого соединения
		}
		go handleClient(conn) // обрабатываем запросы клиента в отдельной го-рутине
	}
}

func handleClient(conn net.Conn) {
	fmt.Println( // вывод на сервер, что замечено новое подключение и время подключения
		"received new connection: ", conn.RemoteAddr(),
		" time: ", time.Now().Format("02-01-2006 15:04:05"))
	defer conn.Close()               // закрываем сокет при выходе из функции
	buf := make([]byte, 1024)        // буфер для чтения клиентских данных
	Vibor(conn, buf)                 // функция для работы с субд
	conn.Write([]byte("\nВыход...")) // пишем в сокет
	fmt.Println(
		"connection closed: ", conn.RemoteAddr(),
		" time: ", time.Now().Format("02-01-2006 15:04:05"))
	fmt.Println()
}
