package main

import (
	"fmt"
	"runtime"
	"time"
	"urlShortener/handler"
	mapstruct "urlShortener/mapStruct"

	"github.com/gin-gonic/gin"
)

const port = "8080"

func main() {
	runtime.GOMAXPROCS(20)

	mapS := mapstruct.New("")
	mapS.Load()

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			mapS.Save()
			fmt.Println("File saved...")
		}
	}()

	router := gin.Default()

	router.GET("/hello", handler.Hello(mapS))

	router.POST("/shorten", handler.Shorten(mapS))

	router.GET("/link/:id", handler.Link(mapS))

	router.Run(":" + port)

}
