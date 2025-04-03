package models

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type CON struct {
	db *sql.DB
}

func Repository(db *sql.DB) *CON {
	return &CON{db: db}
}

type Users struct {
	Id        int
	FirstName string
	LastName  string
}

func Connection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/ecomm_api")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (con *CON) GetAllUsers(db *sql.DB) []Users {
	rows, err := con.db.Query("SELECT * FROM user_api")
	if err != nil {
		log.Fatal("error while fetching", err)

	}
	defer rows.Close()
	var users []Users
	for rows.Next() {
		var user Users
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName)
		if err != nil {
			log.Fatal("error while extracting", err)

		}
		users = append(users, user)
	}
	return users
}
func (con *CON) CreateUser(first_name string, last_name string, db *sql.DB) (string, error) {
	_, err := con.db.Exec("INSERT INTO user_api(first_name,last_name) VALUES(?,?)", first_name, last_name)
	if err != nil {
		return "", err
	}
	return "datainserted", nil
}
