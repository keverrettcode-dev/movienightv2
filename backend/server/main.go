package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	controller "github.com/keverrettcode-dev/movienightv2/backend/server/controllers"
)

func main() {
	//This is the main function

	router:=gin.Default()

	router.GET("/slayer", func(c *gin.Context) {
		c.String(200, "Hello, from the backend, Sr. Developer Everrett!")
	})

	router.GET("/movies", controller.GetMovies())

	if err:=router.Run(":8080"); err != nil{
		fmt.Println("Failed to start server sir!", err)
	}
}