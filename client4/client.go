package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Неверный ввод - Usage: %s host:port ", os.Args[0]) // вывод ошибки если введено больше 2 аргументов
		os.Exit(1)
	}
	serv := os.Args[1]                 // берем адрес сервера из аргументов командной строки
	conn, err := net.Dial("tcp", serv) // открываем TCP-соединение к серверу
	if err != nil {
		fmt.Println(err) // вывод ошибки при открытии TCP-соединения к серверу
		return
	}
	defer conn.Close()
	go copyTo(os.Stdout, conn) // читаем из сокета в stdout
	copyTo(conn, os.Stdin)     // пишем в сокет из stdin
}

func copyTo(dst io.Writer, src io.Reader) {
	io.Copy(dst, src)
}
