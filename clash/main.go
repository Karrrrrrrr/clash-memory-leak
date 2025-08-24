package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"runtime"
	"time"
)

type WrapConn struct {
	net.Conn
	network string
	addr    string
}

func newClient(tcpIn chan<- WrapConn) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				if network != "tcp" && network != "tcp4" && network != "tcp6" {
					return nil, errors.New("unsupported network " + network)
				}

				left, right := net.Pipe()

				tcpIn <- WrapConn{
					Conn:    right,
					network: network,
					addr:    addr,
				}
				return left, nil

			},
		},
	}
	return client
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
			log.Println("NumGoroutine:", runtime.NumGoroutine())
			time.Sleep(time.Second)
		}
	}()
	var tcpIn = make(chan WrapConn, 100)

	client := newClient(tcpIn)

	go func() {
		for conn := range tcpIn {
			conn := conn
			go func() {
				defer conn.Close()
				remoteConn, err := net.Dial("tcp", conn.addr)

				if err != nil {
					log.Println("dial err", err)
					return
				}
				var ch = make(chan error, 1)
				go func() {
					_, err := io.Copy(conn.Conn, remoteConn)
					ch <- err
					//
					remoteConn.SetReadDeadline(time.Now())
					conn.Conn.SetReadDeadline(time.Now())
					log.Println("goroutine done!")
				}()

				_, err = io.Copy(remoteConn, conn.Conn)
				log.Println("out done!")

				remoteConn.SetReadDeadline(time.Now())
				conn.Conn.SetReadDeadline(time.Now())

				if err != nil {
					log.Println("io.Copy", err)
				}
				err = <-ch
				if err != nil {
					log.Println("goroutine io.Copy", err)
				}

			}()
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
			defer cConn.Close()
			bufReader := bufio.NewReader(cConn)
			//bufWriter := bufio.NewWriter(cConn)
			//defer bufWriter.Flush()

			request, err := http.ReadRequest(bufReader)
			if err != nil {
				log.Println("http.readRequest", err)
				return
			}
			request.RemoteAddr = cConn.RemoteAddr().String()
			request.RequestURI = ""
			resp, err := client.Do(request)
			if err != nil {
				log.Println("client.do", err)
				return
			}

			//defer resp.Body.Close()
			var mw = io.MultiWriter(
				//os.Stdout, // for debug
				clientConn,
			)
			err = resp.Write(mw)
			if err != nil {
				//log.Println("resp.write", err)
				return
			}
		}(clientConn)
	}
}
