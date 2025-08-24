package main

import (
	"fmt"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	go func() {
		http.ListenAndServe(":6666", nil)
	}()

	var app = gin.Default()
	app.GET("/v1/health/service/card-service", func(c *gin.Context) {
		x := rand.Int() % 2
		//if x == 1 {
		time.Sleep(time.Second * 1)
		//}
		c.JSON(200, x)
	})

	go func() {
		app.Run()
	}()
	for {
		fmt.Println(runtime.NumGoroutine())
		time.Sleep(time.Second / 2)
	}

}
