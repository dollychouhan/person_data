package main

import (
	"database/sql"
	"log"
)

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
