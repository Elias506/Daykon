package main

import (
	"testing"
	"time"
)

func TestGET(t *testing.T) {
	daykon := &DaykonType{}
	daykon.db = make(map[string]*Elem)
	e1 := &Elem{
		value: []byte("1"),
	}
	daykon.db["key1"] = e1
	e2 := &Elem{
		value: []byte("str2"),
	}
	daykon.db["key2"] = e2

	b, err := daykon.get([][]byte{[]byte("GET"), []byte("key1")})
	if !(err == nil && string(b) == "1") {
		t.Errorf("GetError: %v", b)
	}
	b2, err2 := daykon.get([][]byte{[]byte("GET"), []byte("key2")})
	if !(err2 == nil && string(b2) == "str2") {
		t.Errorf("GetError: %v", b2)
	}
	_, err3 := daykon.get([][]byte{[]byte("GET"), []byte("key2"), []byte("trash")})
	if err3 == nil {
		t.Errorf("GetError 3")
	}
	b4, err4 := daykon.get([][]byte{[]byte("GET"), []byte("key3")})
	if !(b4 == nil && err4 == nil) {
		t.Errorf("GetError 4")
	}
}

func TestSET(t *testing.T) {
	daykon := &DaykonType{}
	daykon.db = make(map[string]*Elem)
	b1, err1 := daykon.set([][]byte{[]byte("SET"), []byte("key1"), []byte("1")})
	if !(err1 == nil && string(b1) == "OK") {
		t.Errorf("SetError 1")
	}
	if _, err2 := daykon.db["key1"]; !err2 {
		t.Errorf("SetError 2")
	}
	err3 := (string(daykon.db["key1"].value) == "1")
	if err3 != true {
		t.Errorf("SetError 3")
	}
	b4, err4 := daykon.set([][]byte{[]byte("SET"), []byte("key2")})
	if !(b4 == nil && err4 != nil) {
		t.Errorf("SetError 4")
	}
	_, _ = daykon.set([][]byte{[]byte("SET"), []byte("key3"), []byte("10"), []byte("1s")})
	time.Sleep(1100 * time.Millisecond)
	if _, err5 := daykon.db["key3"]; err5 {
		t.Errorf("SetError 5")
	}
}

func TestDEL(t *testing.T) {
	daykon := &DaykonType{}
	daykon.db = make(map[string]*Elem)
	e1 := &Elem{
		value: []byte("1"),
	}
	daykon.db["key1"] = e1
	e2 := &Elem{
		value: []byte("str2"),
	}
	daykon.db["key2"] = e2
	e3 := &Elem{
		value: []byte("str2"),
	}
	daykon.db["key3"] = e3

	b1 := daykon.del([][]byte{[]byte("DEL"), []byte("key1")})
	if string(b1) != "1" {
		t.Errorf("DelError 1")
	}
	if _, err2 := daykon.db["key1"]; err2 {
		t.Errorf("DelError 2")
	}

	b2 := daykon.del([][]byte{[]byte("DEL"), []byte("key2"), []byte("key3")})
	if string(b2) != "2" {
		t.Errorf("DelError 3")
	}
	if _, err3 := daykon.db["key2"]; err3 {
		t.Errorf("DelError 3")
	}
	if _, err4 := daykon.db["key4"]; err4 {
		t.Errorf("DelError 4")
	}
}

func TestKEYS(t *testing.T) {
	daykon := &DaykonType{}
	daykon.db = make(map[string]*Elem)
	e1 := &Elem{
		value: []byte("1"),
	}
	daykon.db["name"] = e1
	e2 := &Elem{
		value: []byte("2"),
	}
	daykon.db["firstname"] = e2
	e3 := &Elem{
		value: []byte("3"),
	}
	daykon.db["age"] = e3
	e4 := &Elem{
		value: []byte("4"),
	}
	daykon.db["ag6"] = e4
	_, err0 := daykon.keys([][]byte{[]byte("KEYS")})
	if err0 == nil {
		t.Errorf("KeysError 0")
	}
	b1, err1 := daykon.keys([][]byte{[]byte("KEYS"), []byte(".*name.*")})
	if !(err1 == nil && (string(b1) == `1) name
2) firstname` || string(b1) == `1) firstname
2) name`)) {
		t.Errorf("KeysError 1")
	}
	b2, err2 := daykon.keys([][]byte{[]byte("KEYS"), []byte(`ag\d`)})
	if !(err2 == nil && string(b2) == `1) ag6`) {
		t.Errorf("KeysError 2")
	}
}
