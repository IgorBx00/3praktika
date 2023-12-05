package main

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type Node struct {
	data string
	next *Node
}

type Stack struct { // РЕАЛИЗОВАНО
	head *Node
	mut  sync.Mutex
}

func (stack *Stack) push(val string) { // РАБОТАЕТ
	stack.mut.Lock()
	new_Node := &Node{data: val, next: stack.head}
	stack.head = new_Node
	stack.mut.Unlock()
}

func (stack *Stack) pop() string { // РАБОТАЕТ
	stack.mut.Lock()
	if stack.head == nil {
		stack.mut.Unlock()
		return "Stack is Empty"
	}
	deleted := stack.head.data
	stack.head = stack.head.next
	stack.mut.Unlock()
	return deleted
}

type Queue struct { // РЕАЛИЗОВАНО
	head *Node
	tail *Node
	mut  sync.Mutex
}

func (queue *Queue) Enqueue(val string) { // РАБОТАЕТ
	queue.mut.Lock()
	new_Node := &Node{data: val, next: nil}
	if queue.head == nil {
		queue.head = new_Node
		queue.tail = new_Node
	} else {
		queue.tail.next = new_Node
		queue.tail = new_Node
	}
	queue.mut.Unlock()
}

func (queue *Queue) Dequeue() string { // РАБОТАЕТ
	queue.mut.Lock()
	if queue.head == nil {
		queue.mut.Unlock()
		return "Queue is Empty"
	}
	deleted := queue.head.data
	queue.head = queue.head.next
	if queue.head == nil {
		queue.tail = queue.head
	}
	queue.mut.Unlock()
	return deleted
}

type Pair struct {
	key  string
	url  string
	time string
}

type HashMap struct {
	table [512]*Pair
	mut   sync.Mutex
}

func HashTable(key string, size int) (int, error) {
	if len(key) == 0 {
		return 0, errors.New("keySize 0")
	}
	hashSum := 0
	for i := 0; i < len(key); i++ {
		hashSum += int(key[i])
	}
	return hashSum % size, nil
}

func (hmap *HashMap) insert(key string, value string) string {
	hmap.mut.Lock()
	p := &Pair{key: key, url: value, time: time.Now().Format("02-01-2006 15:04:05")}
	hash, err := HashTable(key, len(hmap.table))
	if err != nil {
		hmap.mut.Unlock()
		return err.Error()
	}
	if hmap.table[hash] == nil {
		hmap.table[hash] = p
		hmap.mut.Unlock()
		return "Добавлено\n"
	} else if hmap.table[hash].key == p.key {
		hmap.table[hash] = p
		hmap.mut.Unlock()
		return "ключ уже существует, значение заменено\n"
	}
	for i := (hash + 1) % len(hmap.table); i != hash; i = (i + 1) % len(hmap.table) {
		if hmap.table[i] == nil {
			hmap.table[i] = p
			hmap.mut.Unlock()
			return "Добавлено\n"
		}
	}
	hmap.mut.Unlock()
	return "full\n"
}

func (hmap *HashMap) get(key string) (string, string) {
	hmap.mut.Lock()
	hash, err := HashTable(key, len(hmap.table))
	if err != nil {
		hmap.mut.Unlock()
		return err.Error(), ""
	}

	for i := hash; hmap.table[i] != nil; i = (i + 1) % len(hmap.table) {
		if hmap.table[i].key == key {
			hmap.mut.Unlock()
			return hmap.table[i].url, hmap.table[i].time
		}
	}
	hmap.mut.Unlock()
	return "not found\n", ""
}

func (hmap *HashMap) del(key string) string {
	hmap.mut.Lock()
	hash, err := HashTable(key, len(hmap.table))
	if err != nil {
		hmap.mut.Unlock()
		return err.Error()
	}

	for i := hash; hmap.table[i] != nil; i = (i + 1) % len(hmap.table) {
		if hmap.table[i].key == key {
			hmap.table[i] = nil
			hmap.mut.Unlock()
			return ""
		}
	}
	hmap.mut.Unlock()
	return "not found\n"
}

type Set struct {
	key string
}
type SetMap struct {
	table [512]*Set
	mut   sync.Mutex
}

func SetHash(key string, size int) (int, error) {
	if len(key) == 0 {
		return 0, errors.New("keySize 0")
	}
	hashSum := 0
	for i := 0; i < len(key); i++ {
		hashSum += int(key[i])
	}
	return hashSum % size, nil
}

func (smap *SetMap) add(key string) string {
	smap.mut.Lock()
	p := &Set{key: key}
	hash, err := SetHash(key, len(smap.table))
	if err != nil {
		smap.mut.Unlock()
		return err.Error()
	}
	if smap.table[hash] == nil {
		smap.table[hash] = p
		smap.mut.Unlock()
		return ""
	} else if smap.table[hash].key == p.key {
		smap.table[hash] = p
		smap.mut.Unlock()
		return "Элемент уже есть\n"
	}
	for i := (hash + 1) % len(smap.table); i != hash; i = (i + 1) % len(smap.table) {
		if smap.table[i] == nil {
			smap.table[i] = p
			smap.mut.Unlock()
			return ""
		}
	}
	smap.mut.Unlock()
	return "full\n"
}

func (smap *SetMap) sismem(key string) string {
	smap.mut.Lock()
	hash, err := SetHash(key, len(smap.table))
	if err != nil {
		smap.mut.Unlock()
		return err.Error()
	}
	for i := hash; smap.table[i] != nil; i = (i + 1) % len(smap.table) {
		if smap.table[i].key == key {
			smap.mut.Unlock()
			return smap.table[i].key
		}
	}
	smap.mut.Unlock()
	return "not found\n"
}

func (smap *SetMap) rem(key string) string {
	smap.mut.Lock()
	hash, err := SetHash(key, len(smap.table))
	if err != nil {
		smap.mut.Unlock()
		return err.Error()
	}

	for i := hash; smap.table[i] != nil; i = (i + 1) % len(smap.table) {
		if smap.table[i].key == key {
			smap.table[i] = nil
			smap.mut.Unlock()
			return ""
		}
	}
	smap.mut.Unlock()
	return "not found\n"
}

func (hmap *HashMap) checkUrl(url string) string {
	for i := 0; i < len(hashtable.table); i++ {
		if hashtable.table[i] != nil && url == hashtable.table[i].url {
			return hashtable.table[i].key
		}
	}
	return "false"
}

// СТЕК
var stack = &Stack{}

// ОЧЕРЕДЬ
var queue = &Queue{}

// ХЕШ-ТАБЛИЦА
var hashtable = &HashMap{}

// МНОЖЕСТВО
var set = &SetMap{}

func Vibor(conn net.Conn, buf []byte) {
	conn.Write([]byte("Command: ")) // пишем в сокет
	readLen, err := conn.Read(buf)  // читаем из сокета
	if err != nil {
		fmt.Println(err)
		return
	}
	var s string
	i := 0
	i2 := 0
	for ; string(buf[i]) != " " && i < readLen-2; i++ {
	}
	s = string(buf[0:i])
	i++
	i2 = i
	if s == "SPOP" || s == "QPOP" {
	} else if i == readLen-2 || i-1 == readLen-2 {
		conn.Write([]byte("incorrectly command\n" + s)) // пишем в сокет
		return
	}
	if s == "SADD" {
		conn.Write([]byte("Добавление элемента во множество\n")) // пишем в сокет
		for ; string(buf[i]) != " " && i < readLen-2; i++ {
		}
		ad := string(buf[i2:i])
		conn.Write([]byte(set.add(ad)))
		fmt.Println(
			"ip: ", conn.RemoteAddr(),
			" command: ", s,
			" Value: ", ad,
			" time: ", time.Now().Format("02-01-2006 15:04:05"))
		fmt.Println()
	} else if s == "SREM" {
		conn.Write([]byte("Удаление элемента из множества\n")) // пишем в сокет
		for ; string(buf[i]) != " " && i < readLen-2; i++ {
		}
		ad := string(buf[i2:i])
		conn.Write([]byte(set.rem(ad)))
		fmt.Println(
			"ip: ", conn.RemoteAddr(),
			" command: ", s,
			" Value: ", ad,
			" time: ", time.Now().Format("02-01-2006 15:04:05"))
		fmt.Println()
	} else if s == "SISMEMBER" {
		conn.Write([]byte("Поиск элемента во множестве: ")) // пишем в сокет
		for ; string(buf[i]) != " " && i < readLen-2; i++ {
		}
		ad := string(buf[i2:i])
		conn.Write([]byte(set.sismem(ad) + "\n"))
		fmt.Println(
			"ip: ", conn.RemoteAddr(),
			" command: ", s,
			" Value: ", ad,
			" time: ", time.Now().Format("02-01-2006 15:04:05"))
		fmt.Println()
		// СТЕК
	} else if s == "SPUSH" {
		conn.Write([]byte("Добавление элемента в стек\n")) // пишем в сокет
		for ; string(buf[i]) != " " && i < readLen-2; i++ {
		}
		ad := string(buf[i2:i])
		stack.push(ad)
		fmt.Println(
			"ip: ", conn.RemoteAddr(),
			" command: ", s,
			" Value: ", ad,
			" time: ", time.Now().Format("02-01-2006 15:04:05"))
		fmt.Println()
	} else if s == "SPOP" {
		conn.Write([]byte("Удаление элемента из стека \n")) // пишем в сокет
		conn.Write([]byte(stack.pop() + "\n"))              // пишем в сокет
		fmt.Println(
			"ip: ", conn.RemoteAddr(),
			" command: ", s,
			" time: ", time.Now().Format("02-01-2006 15:04:05"))
		fmt.Println()
		// ОЧЕРЕДЬ
	} else if s == "QPUSH" {
		conn.Write([]byte("Добавление элемента в очередь: ")) // пишем в сокет
		for ; string(buf[i]) != " " && i < readLen-2; i++ {
		}
		ad := string(buf[i2:i])
		queue.Enqueue(ad)
		fmt.Println(
			"ip: ", conn.RemoteAddr(),
			" command: ", s,
			" Value: ", ad,
			" time: ", time.Now().Format("02-01-2006 15:04:05"))
		fmt.Println()
	} else if s == "QPOP" {
		conn.Write([]byte("Удаление элемента из очереди ")) // пишем в сокет
		conn.Write([]byte(queue.Dequeue() + "\n"))
		fmt.Println(
			"ip: ", conn.RemoteAddr(),
			" command: ", s,
			" time: ", time.Now().Format("02-01-2006 15:04:05"))
		fmt.Println()
		// ХЕШ-ТАБЛИЦА
	} else if s == "HSET" {
		conn.Write([]byte("Добавление ключа и значения в хеш-таблицу\n")) // пишем в сокет                                              // читаем из сокета
		for ; string(buf[i]) != " " && i < readLen-2; i++ {
		}
		key := string(buf[i2:i])
		i++
		i2 = i
		if i == readLen-2 || i-1 == readLen-2 {
			conn.Write([]byte("incorrectly command\n")) // пишем в сокет
			return
		}
		for ; string(buf[i]) != " " && i < readLen-2; i++ {
		}
		ad := string(buf[i2:i])
		conn.Write([]byte(hashtable.insert(key, ad)))
		fmt.Println(
			"ip: ", conn.RemoteAddr(),
			" command: ", s,
			" Value: ", ad,
			" Key: ", key,
			" time: ", time.Now().Format("02-01-2006 15:04:05"))
		fmt.Println()
	} else if s == "HDEL" {
		conn.Write([]byte("Удаление ключа из хеш-таблицы\n")) // пишем в сокет
		for ; string(buf[i]) != " " && i < readLen-2; i++ {
		}
		key := string(buf[i2:i])
		conn.Write([]byte(hashtable.del(key)))
		fmt.Println(
			"ip: ", conn.RemoteAddr(),
			" command: ", s,
			" Key: ", key,
			" time: ", time.Now().Format("02-01-2006 15:04:05"))
		fmt.Println()
	} else if s == "HGET" {
		conn.Write([]byte("Поиск элемента в хеш-таблице: ")) // пишем в сокет
		for ; string(buf[i]) != " " && i < readLen-2; i++ {
		}
		key := string(buf[i2:i])
		bb, _ := hashtable.get(key)
		conn.Write([]byte(bb + "\n"))
		fmt.Println(
			"ip: ", conn.RemoteAddr(),
			" command: ", s,
			" Key: ", key,
			" time: ", time.Now().Format("02-01-2006 15:04:05"))
		fmt.Println()
	} else if s == "HURL" {
		conn.Write([]byte("Поиск элемента в хеш-таблице\n")) // пишем в сокет
		for ; string(buf[i]) != " " && i < readLen-2; i++ {
		}
		url := string(buf[i2:i])
		bb := hashtable.checkUrl(url)
		conn.Write([]byte(bb + "\n"))
		fmt.Println(
			"ip: ", conn.RemoteAddr(),
			" command: ", s,
			" Key: ", url,
			" time: ", time.Now().Format("02-01-2006 15:04:05"))
		fmt.Println()
	} else {
		conn.Write([]byte("incorrectly command\n")) // пишем в сокет
	}
	if err != nil {
		return
	}
	return
}
