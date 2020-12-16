package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

//DaykonType ...
type DaykonType struct {
	db map[string]*Elem
	mu sync.Mutex
}

//Elem ...
type Elem struct {
	value         []byte
	timeStart     time.Time
	timerDuration time.Duration
	timer         *time.Timer
}

func (daykon *DaykonType) write() {
	for k, v := range daykon.db {
		fmt.Printf("%s: %s, ", k, string(v.value))
	}
	fmt.Println()
}

func (daykon *DaykonType) get(source [][]byte) ([]byte, error) {
	if len(source) != 2 {
		return nil, fmt.Errorf("Unexpected input")
	}
	daykon.mu.Lock()
	defer daykon.mu.Unlock()
	key := string(source[1])
	e, ok := daykon.db[key]
	if !ok {
		return nil, nil
	}
	var res []byte
	if e.timer == nil {
		res = e.value
	} else {
		t := e.timerDuration - time.Since(e.timeStart)
		res = bytes.Join([][]byte{e.value, []byte(fmt.Sprint(t))}, []byte(" "))
	}
	return res, nil
}

func (daykon *DaykonType) set(source [][]byte) ([]byte, error) {
	if !(len(source) == 3 || len(source) == 4) {
		return nil, fmt.Errorf("Unexpected input")
	}
	daykon.mu.Lock()
	defer daykon.mu.Unlock()
	key := string(source[1])
	v := source[2]
	e := &Elem{
		value: v,
	}
	daykon.db[key] = e
	if len(source) == 4 {
		duration, err := time.ParseDuration(string(source[3]))
		if err != nil {
			return nil, fmt.Errorf("Unexpected input")
		}
		go ttl(daykon.db, key, duration)
	}
	return []byte("OK"), nil
}

func ttl(db map[string]*Elem, key string, d time.Duration) {
	db[key].timeStart = time.Now()
	db[key].timerDuration = d
	db[key].timer = time.NewTimer(d)
	<-db[key].timer.C
	if _, ok := db[key]; ok {
		fmt.Printf("Time to live 0: %s\n", key)
		delete(db, key)
	}

}

func (daykon *DaykonType) del(source [][]byte) []byte {
	daykon.mu.Lock()
	defer daykon.mu.Unlock()
	k := 0
	for i := 1; i < len(source); i++ {
		key := string(source[i])
		_, ok := daykon.db[key]
		if ok {
			k++
			delete(daykon.db, key)
		}
	}
	return []byte(strconv.Itoa(k))
}

func (daykon *DaykonType) keys(source [][]byte) ([]byte, error) {
	if len(source) != 2 {
		return nil, fmt.Errorf("Unexpected input")
	}
	daykon.mu.Lock()
	defer daykon.mu.Unlock()
	pattern := string(source[1])
	var str []string
	k := 0
	for key := range daykon.db {
		match, err := regexp.Match(pattern, []byte(key))
		if err != nil {
			return nil, err
		}
		if match {
			k++
			str = append(str, fmt.Sprintf("%d) %s", k, key))
		}
	}
	if str == nil {
		return nil, nil
	}
	join := strings.Join(str, "\n")
	return []byte(join), nil
}

func (daykon *DaykonType) save(source [][]byte) ([]byte, error) {
	if len(source) != 2 {
		return nil, fmt.Errorf("Unexpected input")
	}
	daykon.mu.Lock()
	defer daykon.mu.Unlock()
	file, err := os.Create(fmt.Sprint(string(source[1]), ".bin"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	for key, value := range daykon.db {
		var s string
		if value.timer == nil {
			s = fmt.Sprintf("%v;%v;%v\n", key, string(value.value), "nil")
		} else {
			s = fmt.Sprintf("%v;%v;%v\n", key, string(value.value), value.timerDuration-time.Since(value.timeStart))
		}
		_, err := file.WriteString(s)
		if err != nil {
			return nil, err
		}
	}
	return []byte("OK"), nil
}

func (daykon *DaykonType) backup(source [][]byte) ([]byte, error) {
	if len(source) != 2 {
		return nil, fmt.Errorf("Unexpected input")
	}
	file, err := os.Open(fmt.Sprint(string(source[1]), ".bin"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	daykon.mu.Lock()
	defer daykon.mu.Unlock()
	data := make([]byte, 64)

	for k := range daykon.db {
		delete(daykon.db, k)
	}
	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		dataSplit := bytes.Split(data[0:n], []byte("\n"))
		for _, str := range dataSplit[0 : len(dataSplit)-1] {
			strSplit := bytes.Split(str, []byte(";"))
			key := string(strSplit[0])
			fmt.Println("Backup key:", key)
			value := strSplit[1]
			e := &Elem{
				value: value,
			}
			daykon.db[key] = e
			if t := string(strSplit[2]); t != "nil" {
				d, err := time.ParseDuration(t)
				if err != nil {
					return nil, err
				}
				go ttl(daykon.db, key, d)
			}
			fmt.Printf("Backup: %v - %v\n", key, string(value))
		}

	}

	return []byte("OK"), nil
}

func main() {
	daykon := &DaykonType{}
	daykon.db = make(map[string]*Elem)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening...")
	for {
		conn, err := listener.Accept()
		addr := conn.RemoteAddr().String()
		fmt.Printf("User %s connected\n", addr)
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		go handler(conn, daykon, addr)
	}
}

func handler(conn net.Conn, daykon *DaykonType, addr string) {
	defer conn.Close()
	for {
		req := make([]byte, (1024 * 8))
		n, err := conn.Read(req)
		if n == 0 || err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("%s> %s\n", addr, string(req[0:n]))
		source := bytes.Split(req[0:n], []byte(" "))
		var res []byte
		switch string(source[0]) {
		case "GET":
			res, err = daykon.get(source)
		case "SET":
			res, err = daykon.set(source)
		case "DEL":
			res = daykon.del(source)
			err = nil
		case "KEYS":
			res, err = daykon.keys(source)
		case "SAVE":
			res, err = daykon.save(source)
		case "BACKUP":
			res, err = daykon.backup(source)
		default:
			res, err = nil, fmt.Errorf("Unexpected input")
		}
		if err != nil {
			fmt.Println("Error:", err)
			conn.Write([]byte(fmt.Sprintf("Error: %v", err)))
			continue
		}
		if res == nil {
			res = []byte("nil")
		}
		conn.Write(res)
	}
}
