package cluster

import (
	zmq "github.com/pebbe/zmq4"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	s "strings"
)

const (
	BROADCAST = -1
)

type Envelope struct {
	// On the sender side, Pid identifies the receiving peer. If instead, Pid is
	// set to cluster.BROADCAST, the message is sent to all peers. On the receiver side, the
	// Id is always set to the original sender. If the Id is not found, the message is silently dropped
	Pid int

	// An id that globally and uniquely identifies the message, meant for duplicate detection at
	// higher levels. It is opaque to this package.
	MsgId int64

	// the actual message.
	Msg interface{}
}

type Info struct {
	Myid    int
	Peer    [100]int
	Clstr   map[int]*zmq.Socket
	MyClstr *zmq.Socket
	Outbx   chan *Envelope
	Inbx    chan *Envelope
	len     int
}

type Server interface {
	// Id of this server
	Pid() int

	// array of other servers' ids in the same cluster
	Peers() []int

	// the channel to use to send messages to other peers
	// Note that there are no guarantees of message delivery, and messages
	// are silently dropped
	Outbox() chan *Envelope

	// the channel to receive messages from other peers.
	Inbox() chan *Envelope
}

func (c Info) Pid() int {
	return c.Myid
}

func (c Info) Peers() [100]int {
	return c.Peer
}

func (c Info) Outbox() chan *Envelope {
	return c.Outbx
}

func (c Info) Inbox() chan *Envelope {
	return c.Inbx
}

type jsonobject struct {
	Object ObjectType
}

type ObjectType struct {
	Peers []Peer
}

type Peer struct {
	Id  int
	Url string
}

var myInfo Info

func (c Info) Send() {
	for {
		data, _ := <-c.Outbox()
		if data.Pid == -1 {
			for i := 0; i < c.len; i++ {
				c.Clstr[c.Peer[i]].Send(fmt.Sprintf("%d", c.Pid())+":"+data.Msg.(string), zmq.DONTWAIT)
			}
		} else {
			c.Clstr[data.Pid].Send(fmt.Sprintf("%d", c.Pid())+":"+data.Msg.(string), zmq.DONTWAIT)
		}
	}
}

func (c Info) Recieve() {
	for {
		msg, _ := c.MyClstr.Recv(0)
		id1, _ := strconv.Atoi(s.Split(msg, ":")[0])
		msg1 := msg[len(s.Split(msg, ":")[0])+1:]
		c.Inbox() <- &Envelope{Pid: id1, Msg: string(msg1)}
	}
}

func New(myid int, jsonFile string) Info {
	//fmt.Println(myid)
	myInfo.len = 0
	file, e := ioutil.ReadFile(jsonFile)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var jsontype jsonobject
	json.Unmarshal(file, &jsontype)
	myInfo.Myid = myid
	myInfo.Clstr = make(map[int]*zmq.Socket)
	peerInfo := jsontype.Object.Peers
	j := 0
	for i := 0; i < len(peerInfo); i++ {
		if peerInfo[i].Id != myid {
			myInfo.Peer[j] = peerInfo[i].Id
			myInfo.Clstr[peerInfo[i].Id], _ = zmq.NewSocket(zmq.PUSH)
			myInfo.Clstr[peerInfo[i].Id].Connect("tcp://" + peerInfo[i].Url)

			j++
		} else {
			myInfo.MyClstr, _ = zmq.NewSocket(zmq.PULL)
			myInfo.MyClstr.SetIdentity(string(myid))
			myInfo.MyClstr.Bind("tcp://*:" + s.Split(peerInfo[i].Url, ":")[1])
		}
	}
	myInfo.len = j
	myInfo.Inbx = make(chan *Envelope)
	myInfo.Outbx = make(chan *Envelope)

	return myInfo
}
