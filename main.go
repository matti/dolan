package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mxk/go-flowrate/flowrate"
	"github.com/r3labs/sse/v2"
)

func main() {

	if endpoint := os.Getenv("STREAM_FROM"); endpoint != "" {
		go func() {
			client(endpoint)
		}()
	}

	server(os.Getenv("LADDR"))
}

var monitor *flowrate.Monitor

func client(endpoint string) {
	client := sse.NewClient(endpoint)

	monitor = flowrate.New(time.Second*1, time.Second*1)

	go func() {
		for {
			fmt.Println(monitor.Status().InstRate/(1<<20), "mb/s")
			time.Sleep(time.Second * 1)
		}
	}()
	client.Subscribe("messages", func(msg *sse.Event) {
		monitor.Update(len(msg.Data))
	})

}

func server(laddr string) {
	r := gin.Default()
	r.LoadHTMLGlob("./views/**")

	r.GET("/stream", func(c *gin.Context) {
		c.Stream(func(w io.Writer) bool {
			select {
			case <-c.Done():
				return false
			case <-c.Request.Context().Done():
				return false
			default:
				data := make([]byte, 1024)
				rand.Read(data)

				c.SSEvent("data", data)
				return true
			}
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"rate": monitor.Status().CurRate / (1 << 20),
			"from": os.Getenv("STREAM_FROM"),
		})

	})

	r.Run(laddr)
}
