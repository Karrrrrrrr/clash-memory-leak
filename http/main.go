package main

import (
	"context"
	"fmt"
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
	for {
		c++
		func() {
			var ctx, _ = context.WithTimeout(context.Background(), time.Second/4)

			req, err := http.NewRequestWithContext(ctx, "GET", "http://192.168.88.1:8080/v1/health/service/card-service", nil)
			//req, err := http.NewRequestWithContext(ctx, "GET", "http://192.168.88.1:8080/v1/health/service/card-service", nil)

			_ = err

			do, err := client.Do(req)
			if err != nil {
				//panic(err)
				fmt.Println(err)
			}
			if err == nil {
				defer do.Body.Close()
			}
			_ = do
			time.Sleep(time.Second / 3)
		}()

	}

}
