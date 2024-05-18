package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	//Gin router
	router := gin.Default()

	//Database connection
	dsn := "username:password@tcp(127.0.0.1:3306)/cetec"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Error while connecting the MYSQL database: ", err)
		return
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println("Error while verifying the connection to the databse: ", err)
		return
	}

	registerRoutes(router, db)

	//Start the server
	router.Run(":8082")

}
