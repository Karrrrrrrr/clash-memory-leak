package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

func main() {
	c := 0
	go func() {
		// for debug
		for {
			time.Sleep(time.Second)
			fmt.Println(runtime.NumGoroutine())
		}
	}()
	proxyURL, err := url.Parse("http://localhost:8887")
	if err != nil {
		panic(err)
	}
	proxy := http.ProxyURL(proxyURL)

	var client = &http.Client{
		Transport: &http.Transport{Proxy: proxy},
	}
	//http.ProxyFromEnvironment()
	for {
		c++
		func() {
			var ctx = context.Background()
			ctx, _ = context.WithTimeout(ctx, time.Second/10)

			req, err := http.NewRequestWithContext(ctx, "GET", "http://192.168.88.1:8080/v1/health/service/card-service", nil)
			//req, err := http.NewRequestWithContext(ctx, "GET", "http://192.168.88.1:8080/v1/health/service/card-service", nil)

			_ = err
			//net.Pipe()

			do, err := client.Do(req)
			if err != nil {
				//panic(err)
				fmt.Println(err)
			}
			if err == nil {
				all, err := io.ReadAll(do.Body)
				fmt.Println("string(all):", string(all), err)
				defer do.Body.Close()
			}

			_ = do
			//time.Sleep(time.Second / 3)
		}()

	}

}
