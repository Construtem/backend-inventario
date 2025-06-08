package main

import "fmt"
import "github.com/gin-gonic/gin"
import "github.com/joho/godotenv"
import "os"
import "log"
import "db/db.go"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar .env")
	}

	db.Conectar()
	fmt.Println("Conexión a la base de datos exitosa")

  	router := gin.Default()
//	router.Static("/static", "./static")

  	router.GET("/", func(c *gin.Context) {
  	  	c.JSON(200, gin.H{
  	    	"message": "pong",
  	  	})
  	})
  	router.Run() // listen and serve on 0.0.0.0:8080
}