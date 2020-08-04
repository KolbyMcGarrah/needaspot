package interests

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	database "github.com/KolbyMcGarrah/nas/internal/pkg/db/postgres"
	"github.com/KolbyMcGarrah/nas/internal/requests"
	"github.com/KolbyMcGarrah/nas/internal/users"
)

type Interest struct {
	ID           int
	User         *users.User
	AcceptedUser *users.User
	Request      *requests.Request
	Description  string
	Accepted     bool
	CreatedTs    string
	AcceptedTs   string
}

//Save creates a request in the database and returns it's ID.
func (i *Interest) Save() int {
	//Prepare the sql statement
	stmt, err := database.Db.Prepare("INSERT INTO interests(user_id, request_id, description) VALUES($1, $2, $3) RETURNING interest_Id")
	if err != nil {
		log.Fatalf("Error preparing save interest sql: %s", err)
	}
	var ID int
	usr_id, _ := strconv.Atoi(i.User.ID)
	req_id, _ := strconv.Atoi(i.Request.ID)
	err = stmt.QueryRow(usr_id, req_id, i.Description).Scan(&ID)
	if err != nil {
		log.Panicf("Error writing Interest to the database: %s", err)
	}
	return ID
}

//GetInterestByReqID grabs all of the interests from the database associated with a request.
func GetInterestByReqID(r requests.Request) []Interest {
	//prepare sql statement
	s, err := database.Db.Prepare("SELECT i.interest_id, i.description, i.accepted, i.created_ts, i.accepted_ts, " +
		"u.user_id, u.username, u.gender, u.age, u.level, " +
		"a.user_id, a.username, a.gender, a.age, a.level " +
		"FROM interests i " +
		"join users u on u.user_id = i.user_id " +
		"left join users a on a.user_id = i.accepted_user " +
		"WHERE request_id = $1")
	if err != nil {
		log.Fatalf("Error preparing sql for GetInterestsByReqID: %s", err)
	}
	//create the return slice.
	var interests []Interest
	//query for all rows matching req_id.
	defer s.Close()
	rows, err := s.Query(r.ID)
	if err != nil {
		log.Panicf("Error executing GetInterestByReqID: %s", err)
	}

	//create parse variables for the user
	var user_id string
	var username string
	var gender string
	var age int
	var level string

	//create parse variables for the accepted user
	var accid sql.NullString
	var accname sql.NullString
	var accgender sql.NullString
	var accage sql.NullString
	var acclevel sql.NullString

	//create nullStrings for the AcceptedTs
	var acc_ts sql.NullString

	log.Printf("Returned %v", rows)
	for rows.Next() {
		var i Interest
		err = rows.Scan(&i.ID, &i.Description, &i.Accepted, &i.CreatedTs, &acc_ts, &user_id, &username, &gender, &age, &level, &accid, &accname, &accgender, &accage, &acclevel)
		if err != nil {
			log.Panicf("Unable to map interest from database in GetInterestByReqID: %s", err)
		} else {
			//create the interest user
			i.User = &users.User{
				ID:       user_id,
				Username: username,
				Gender:   gender,
				Age:      age,
				Level:    level,
			}
			//Map accepted user if accepted is true
			if i.Accepted {
				//create new user
				var auser users.User
				//validate and parse all NullStrings
				if accid.Valid {
					auser.ID = accid.String
				}
				if accname.Valid {
					auser.Username = accname.String
				}
				if accgender.Valid {
					auser.Username = accgender.String
				}
				if accage.Valid {
					auser.Age, _ = strconv.Atoi(accage.String)
				}
				if acclevel.Valid {
					auser.Level = acclevel.String
				}
				//Assign auser to Accepted user
				i.AcceptedUser = &auser
			}
			if acc_ts.Valid {
				i.AcceptedTs = acc_ts.String
			}
			interests = append(interests, i)
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return interests

}

//GetInterestsByUser returns all interests by user_id
func GetInterestsByUser(u users.User) []Interest {
	//Prepare SQL statement
	s, err := database.Db.Prepare("SELECT i.interest_id, i.description, i.accepted, i.created_ts, i.accepted_ts, " +
		"r.request_id, r.title, r.location, r.workout, " +
		"c.user_id, c.username, c.gender, c.age, c.level, " +
		"a.user_id, a.username, a.gender, a.age, a.level " +
		"FROM interests i " +
		"join requests r on r.request_id = i.request_id " +
		"join users c on c.user_id = r.creator " +
		"left join users a on a.user_id = i.accepted_user " +
		"WHERE i.user_id = $1")
	if err != nil {
		log.Fatalf("Error preparing sql statement for GetInterestsByUser: %s", err)
	}
	//execute the sql
	defer s.Close()
	rows, err := s.Query(u.ID)
	if err != nil {
		log.Panicf("Error returning rows for GetInterestsByUser: %s", err)
	}
	//create return variable
	var interests []Interest

	//Create NullStrings for Accepted User and AcceptedTs
	var u_id sql.NullString
	var username sql.NullString
	var gender sql.NullString
	var age sql.NullString
	var level sql.NullString

	var acc_ts sql.NullString

	//Create an Interest for every row returned
	for rows.Next() {
		//Create Interest
		var i Interest
		//Create Request
		var r requests.Request
		//Create User for Request
		var c users.User
		err = rows.Scan(&i.ID, &i.Description, &i.Accepted, &i.CreatedTs, &acc_ts, &r.ID, &r.Title, &r.Workout, &r.Location,
			&c.ID, &c.Username, &c.Gender, &c.Age, &c.Level, &u_id, &username, &gender, &age, &level)
		//map the user to the request
		r.User = &c
		//map the request to the interest
		i.Request = &r
		//map the user to the interest
		i.User = &u

		//If the interest is accepted, created and map the accepted user.
		if i.Accepted {
			//Create user
			var a users.User
			//Validate NullString and map to user
			if u_id.Valid {
				a.ID = u_id.String
			}
			if username.Valid {
				a.Username = username.String
			}
			if gender.Valid {
				a.Gender = gender.String
			}
			if age.Valid {
				a.Age, _ = strconv.Atoi(age.String)
			}
			if level.Valid {
				a.Level = level.String
			}
		}
		if err != nil {
			log.Panicf("Error mapping interest with ID: %v due to %s", i.ID, err)
		}
		interests = append(interests, i)
	}
	return interests
}

//Accept sets the accept boolean on the interest to true, updates the accepted_ts and sets the accepted user
func (i *Interest) Accept(u users.User) {
	//Prepare the sql statment
	s, err := database.Db.Prepare("UPDATE interests SET accepted=TRUE, accepted_ts=NOW(), accepted_user=$1 WHERE interest_id=$2")
	if err != nil {
		log.Fatalf("Error preparing sql for Accept: %s", err)
	}
	_, err = s.Exec(u.ID, i.ID)
	if err != nil {
		log.Fatalf("Error updating Interest with ID: %v due to %s", i.ID, err)
	}
	i.Accepted = true
	i.AcceptedUser = &u
	i.AcceptedTs = time.Now().String()
}
