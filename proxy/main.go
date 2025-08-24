package main

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"time"
)

func main() {

	listen, err := net.Listen("tcp", ":8887")
	if err != nil {
		panic(err)
	}

	// 用于debug
	go func() {
		for {
			fmt.Println("NumGoroutine:", runtime.NumGoroutine())
			time.Sleep(time.Second)
		}
	}()
	// 这里是核心
	for {
		clientConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go func() {
			remoteConn, err := net.Dial("tcp", "127.0.0.1:8080")
			if err != nil {
				return
			}
			// 核心中的核心 pipe部分
			left, right := net.Pipe()
			var ch = make(chan error)
			go func() {
				_, err := io.Copy(right, remoteConn)
				right.SetReadDeadline(time.Now())
				ch <- err
			}()
			_, _ = io.Copy(left, clientConn)
			left.SetReadDeadline(time.Now())
			<-ch
		}()

	}
}
