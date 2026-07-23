package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	//This is the main function

	router:=gin.Default()
	router.GET("/slayer", func(c *gin.Context) {
		c.String(200, "Hello, from the backend, Sr. Developer Everrett!")
	})

	if err:=router.Run(":8080"); err != nil{
		fmt.Println("Failed to start server sir!", err)
	}
}