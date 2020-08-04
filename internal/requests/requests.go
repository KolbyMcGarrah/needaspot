package requests

import (
	"log"
	"strconv"

	database "github.com/KolbyMcGarrah/nas/internal/pkg/db/postgres"
	"github.com/KolbyMcGarrah/nas/internal/users"
)

type Request struct {
	ID       string      `json:"id"`
	Title    string      `json:"title"`
	Location string      `json:"location"`
	Workout  string      `json:"workout"`
	User     *users.User `json:"user"`
}

// I will need to add user id to this function later
// Save inserts a new request into the database and returns the newly created id
func (request Request) Save() int {
	// prepare the sql query
	statement, err := database.Db.Prepare("INSERT INTO requests(title, location, workout, creator) VALUES($1,$2,$3,$4) RETURNING request_id;")
	if err != nil {
		log.Fatal("SQL statement prep error: ", err)
	}
	// Execute the statement
	var id int
	err = statement.QueryRow(request.Title, request.Location, request.Workout, request.User.ID).Scan(&id)

	if err != nil {
		log.Fatal("Statement execution error: ", err)
	}
	return id
}

func GetReqByID(id int) Request {
	//prepare the sql query
	stmt, err := database.Db.Prepare("SELECT title, location, workout, creator FROM requests WHERE request_id = $1")
	if err != nil {
		log.Fatal("Error preparing GetReqByID: ", err)
	}
	row := stmt.QueryRow(id)
	var request Request
	var user_id int
	err = row.Scan(&request.Title, &request.Location, &request.Workout, &user_id)
	user := users.GetUserById(user_id)
	request.User = &user
	if err != nil {
		log.Fatal("Failed to retrieve request: ", err)
	}
	request.ID = strconv.Itoa(id)
	return request
}

func GetAll() []Request {
	stmt, err := database.Db.Prepare("select r.request_id, r.title, r.workout, r.location, u.user_id, u.username, u.age, u.gender, u.level from requests r inner join users u on r.creator = u.user_id")
	if err != nil {
		log.Fatal("Error preparing GetAll (requests) sql: ", err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("Error executing GetAll (requests) sql: ", err)
	}
	defer rows.Close()
	var requests []Request
	var username string
	var id string
	var age int
	var gender string
	var level string
	for rows.Next() {
		var request Request
		err := rows.Scan(&request.ID, &request.Title, &request.Workout, &request.Location, &id, &username, &age, &gender, &level)
		if err != nil {
			log.Fatal("Error Scanning requests returned from GetAll (requests): ", err)
		}
		request.User = &users.User{
			ID:       id,
			Username: username,
			Age:      age,
			Gender:   gender,
			Level:    level,
		}
		requests = append(requests, request)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return requests
}
