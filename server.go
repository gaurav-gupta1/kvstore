package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	s "strings"
)

var keyval = make(map[string]string)
var mutex = &sync.Mutex{}

func main() {
	service := ":1200"
	tcpaddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpaddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	var buf [512]byte
	n, err := conn.Read(buf[0:])
	checkError(err)
	qry := s.Split(string(buf[0:n])," ")
	
	if qry[0] == "set" {
		mutex.Lock()
		keyval[qry[1]] = qry[2]
		mutex.Unlock()
		_, err := conn.Write([]byte("Key Value Set"))
		checkError(err)
	} else if qry[0] == "get" {
		if val, ok := keyval[qry[1]] ; ok {
			mutex.Lock()
			_, err := conn.Write([]byte(qry[1] + " : " + val))
			mutex.Unlock()
			checkError(err)
		} else {
			_, err := conn.Write([]byte("Key not found in store"))
			checkError(err)
		}
	} else if qry[0] == "delete" {
		if _, ok := keyval[qry[1]] ; ok {
			mutex.Lock()
			delete(keyval,qry[1])
			mutex.Unlock()
			_, err := conn.Write([]byte("Key "+qry[1]+" deleted"))
			checkError(err)
		} else {
			_, err := conn.Write([]byte("Key not found in store"))
			checkError(err)
		}
	} else {
		_, err := conn.Write([]byte("Invalid Command"))
		checkError(err)
	}
	
	if err != nil {
		return
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
