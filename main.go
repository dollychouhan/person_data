package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// struct for getting the person information
type Person struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	City        string `json:"city"`
	State       string `json:"state"`
	Street1     string `json:"street1"`
	Street2     string `json:"street2"`
	ZipCode     string `json:"zip_code"`
}

// struct for inserting new data
type NewPerson struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	City        string `json:"city"`
	State       string `json:"state"`
	Street1     string `json:"street1"`
	Street2     string `json:"street2"`
	ZipCode     string `json:"zip_code"`
}

// CreatePerson function is used to create a new person with data in the database
func CreatePerson(db *sql.DB, newPerson NewPerson) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println("Error while quering the database: ", err)
		return err
	}

	// Insert name into person table
	res, err := tx.Exec("INSERT INTO person (name) VALUES (?)", newPerson.Name)
	if err != nil {
		log.Println("Error while inserting the data into person table")
		tx.Rollback()
		return err
	}
	personID, err := res.LastInsertId()
	if err != nil {
		log.Println("Error while getting the last inserted id of person")
		tx.Rollback()
		return err
	}

	// Insert phone_number into phone table
	_, err = tx.Exec("INSERT INTO phone (number, person_id) VALUES (?, ?)", newPerson.PhoneNumber, personID)
	if err != nil {
		log.Println("Error while inserting the phone number in phone table")
		tx.Rollback()
		return err
	}

	// Insert address data into address table
	res, err = tx.Exec("INSERT INTO address (city, state, street1, street2, zip_code) VALUES (?, ?, ?, ?, ?)",
		newPerson.City, newPerson.State, newPerson.Street1, newPerson.Street2, newPerson.ZipCode)
	if err != nil {
		tx.Rollback()
		return err
	}
	addressID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert joining of person with address into address_join table
	_, err = tx.Exec("INSERT INTO address_join (person_id, address_id) VALUES (?, ?)", personID, addressID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// GetPersonInfo function is used to fetch the data of person from database
func GetPersonInfo(db *sql.DB, personId string) (Person, error) {
	var person Person

	err := db.QueryRow(`SELECT p.name, ph.number, a.city, a.state, a.street1, a.street2, a.zip_code
	FROM person p
	JOIN phone ph ON p.id = ph.person_id
	JOIN address_join aj ON p.id = aj.person_id
	JOIN address a ON aj.address_id = a.id
	WHERE p.id = ? `, personId).Scan(&person.Name, &person.PhoneNumber, &person.City, &person.State, &person.Street1, &person.Street2, &person.ZipCode)

	if err != nil {
		log.Println("Error while fetching the data from database")
		return person, err
	}

	return person, nil

}

func registerRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/person/:person_id/info", func(c *gin.Context) {
		personId := c.Param("person_id")
		person, err := GetPersonInfo(db, personId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, person)
	})

	router.POST("/person/create", func(c *gin.Context) {
		var newPerson NewPerson
		if err := c.ShouldBindJSON(&newPerson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := CreatePerson(db, newPerson)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Person created successfully"})
	})
}

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
