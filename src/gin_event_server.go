package main

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
)

type person struct {
	name string
	age  int
}

func main() {
	ps := person{}
	ps.name = "gildong"
	ps.age = 20

	//var sData []data.Group
	var psData []person
	psData = append(psData, ps)
	psData = append(psData, ps)

	gin.DefaultWriter = colorable.NewColorableStderr()
	r := gin.Default()
	r.GET("/stream", func(c *gin.Context) {
		chanStream := make(chan int, 10)
		go func() {
			defer close(chanStream)
			//for i := 0; i < 5; i++ {
			i := 0
			for {
				chanStream <- i
				i++
				time.Sleep(time.Second * 1)
			}
		}()
		c.Stream(func(w io.Writer) bool {
			if msg, ok := <-chanStream; ok {
				log.Println(msg)
				log.Println(psData)
				//c.SSEvent("message", msg)
				//c.SSEvent("message", psData)
				// c.Render(-1, sse.Event{
				// 	Event: "hola",
				// })
				//c.Render(200, render.JSON{psData})

				user := &person{name: "Frank", age: 20}
				b, err := json.Marshal(user)
				if err != nil {
					log.Println(err)
				}
				log.Println(string(b))

				c.Render(-1, sse.Event{
					Event: "message",
					Data:  string(b),
				})

				return true
			}
			return false
		})
	})
	r.Use(static.Serve("/", static.LocalFile("./public", true)))
	r.Run(":7000")
}
