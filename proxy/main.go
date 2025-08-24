package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"runtime"
	"time"
)

// 带超时的Reader包装器：确保Read操作不会无限阻塞
type timeoutReader struct {
	net.Conn
	timeout time.Duration // 每次读取的超时时间
}

// Read 每次读取前设置超时，确保阻塞不超过timeout
func (t *timeoutReader) Read(p []byte) (int, error) {
	// 每次读取前更新超时时间（从当前时间开始计算）
	if err := t.Conn.SetReadDeadline(time.Now().Add(t.timeout)); err != nil {
		return 0, err
	}
	return t.Conn.Read(p)
}

func main() {
	listen, err := net.Listen("tcp", ":8887")
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	// 监控goroutine数量
	go func() {
		for {
			fmt.Println("NumGoroutine:", runtime.NumGoroutine())
			time.Sleep(time.Second)
		}
	}()

	// 核心逻辑
	for {
		clientConn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go func(cConn net.Conn) {
			// 确保客户端连接最终关闭（无论是否发生错误）
			defer cConn.Close()

			// 连接到目标服务
			remoteConn, err := net.Dial("tcp", "127.0.0.1:8080")
			if err != nil {
				return
			}
			defer remoteConn.Close() // 确保远程连接关闭
			defer clientConn.Close() // 确保远程连接关闭

			left, right := net.Pipe()
			ch := make(chan error, 1)
			go func() {
				_, err := io.Copy(right, remoteConn)
				right.SetReadDeadline(time.Now())

				ch <- err // 无论成功/失败，都通知主goroutine
			}()
			go func() {
				_, _ = io.Copy(remoteConn, right)
			}()
			go func() {
				_, _ = io.Copy(clientConn, left)
			}()
			_, err = io.Copy(left, clientConn)
			left.SetReadDeadline(time.Now())
			//left.Close()
			//right.Close()
			fmt.Println("left finish")
			<-ch
		}(clientConn)
	}
}
