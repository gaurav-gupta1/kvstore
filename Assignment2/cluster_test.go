package cluster

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
	"time"
)

const (
	SERVER = 3
)
const (
	msg_recieved = 100*4 + 100*4 + 4
)

func TestSend(t *testing.T) {
	var id [3]int
	for i := 0; i < 3; i++ {
		id[i] = i + 1
	}

	done := make(chan string, 5)

	for _, myid1 := range id {
		go func(myid int) {
			count := 0
			server := New(myid, "./config.json")
			go server.Send()
			go server.Recieve()

			go func(myid int, count int) {
				for {
					select {
					case <-server.Inbox():
						count = count + 1
					case <-time.After(10 * time.Second):
						done <- "1"
						time.After(5 * time.Second)
					}
				}
			}(myid, count)

			for i := 0; i < 100; i++ {
				sender_pid := "-1"
				id, err := strconv.Atoi(sender_pid)
				msg := fmt.Sprintf("%d", myid)
				if err == nil {
					server.Outbox() <- &Envelope{Pid: id, Msg: msg}
				}
			}
			time.Sleep(time.Second)

			for i := 0; i < 100; i++ {
				for _, myid = range id {
					sender_pid := fmt.Sprintf("%d", myid)
					id, err := strconv.Atoi(sender_pid)
					msg := fmt.Sprintf("%d", myid)
					if err == nil {
						server.Outbox() <- &Envelope{Pid: id, Msg: msg}
					}
				}
			}

			time.Sleep(time.Second)

			file, e := ioutil.ReadFile("./input.txt")
			if e != nil {
				fmt.Printf("File error: %v\n", e)
			}
			for _, myid = range id {
				sender_pid := fmt.Sprintf("%d", myid)
				id, err := strconv.Atoi(sender_pid)
				msg := string(file)
				if err == nil {
					server.Outbox() <- &Envelope{Pid: id, Msg: msg}
				}
			}
		}(myid1)
	}
	count := 0
	for _ = range done {
		count = count + 1
		if count == 3 {
			break
		}
	}
}
