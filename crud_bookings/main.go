package main

import (
	"fmt"
	"log"

	"github.com/somyaranjan99/crud/models"
)

func main() {
	con, err := models.Connection()
	if err != nil {
		log.Fatal("error while initialisize", err)
	}
	respo := models.Repository(con)

	res, err := respo.CreateUser("Radhe", "das", con)
	if err != nil {
		fmt.Println("error while initialisize", err)
	}
	fmt.Println(res)
	data := respo.GetAllUsers(con)
	fmt.Println(data)
}
