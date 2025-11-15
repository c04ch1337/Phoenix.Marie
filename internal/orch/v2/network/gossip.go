package network

import (
	"fmt"
	"log"
	"net"
	"time"
)

var (
	GossipPort      = "9001"
	IsServerRunning bool
)

func StartGossipServer(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("Gossip server error: %v\n", err)
		return
	}
	IsServerRunning = true
	log.Printf("Gossip server up: %s\n", addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}
	msg := string(buf[:n])
	log.Printf("Gossip received: %s\n", msg)
}

func Broadcast(msg string) {
	if !IsServerRunning {
		return
	}
	conn, err := net.DialTimeout("tcp", "localhost:"+GossipPort, 2*time.Second)
	if err != nil {
		log.Printf("Gossip broadcast failed: %v\n", err)
		return
	}
	defer conn.Close()
	conn.Write([]byte(msg))
}

func Heartbeat(agentID string) {
	msg := fmt.Sprintf("HEARTBEAT:%s:%d", agentID, time.Now().Unix())
	Broadcast(msg)
}
