package comments

import (
	"log"
	"strconv"

	"github.com/KolbyMcGarrah/nas/internal/requests"
	"github.com/KolbyMcGarrah/nas/internal/users"

	database "github.com/KolbyMcGarrah/nas/internal/pkg/db/postgres"
)

type Comment struct {
	ID   string           `json:"id"`
	User users.User       `json:"user"`
	Req  requests.Request `json:"req"`
	Text string           `json:"text"`
}

//Save stores the new comment in the database.
func (com Comment) Save() int {
	//Prepare the sql statement
	stmt, err := database.Db.Prepare("INSERT INTO comments(commenter, req, comment) VALUES($1, $2, $3) RETURNING comment_id")
	if err != nil {
		log.Fatal("Error preparing sql to save comment: ", err)
	}

	//create variable to store id of new comment
	var ID int

	usr_id, err := strconv.Atoi(com.User.ID)
	//save comment to database and return id in the
	err = stmt.QueryRow(usr_id, com.Req.ID, com.Text).Scan(&ID)
	if err != nil {
		log.Fatal("Error saving comment to database: ", err)
	}
	return ID
}

//RetrieveRequestComments gets all of the comments belonging to a request.
func RetrieveRequestComments(r requests.Request) []Comment {
	//Prepare sql
	stmt, err := database.Db.Prepare("SELECT c.comment_id, c.comment, u.user_id, u.username, u.level, u.gender, u.age FROM comments c JOIN users u ON u.user_id = c.commenter WHERE req=$1")
	if err != nil {
		log.Fatal("Error preparing sql statment: ", err)
	}
	defer stmt.Close()
	//retrieve all results
	rows, err := stmt.Query(r.ID)
	if err != nil {
		log.Fatal("Error executing RetrieveComments sql: ", err)
	}
	defer rows.Close()
	var comments []Comment
	//initialize variables to construct comments
	var ID int
	var text string
	var user_id string
	var level string
	var gender string
	var age int
	var username string
	//Go through each row and append a comment to the return value
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&ID, &text, &user_id, &username, &level, &gender, &age)
		if err != nil {
			log.Fatal("Error with return value types on GetComments: ", err)
		}
		comment.User = users.User{
			ID:       user_id,
			Username: username,
			Age:      age,
			Gender:   gender,
			Level:    level,
		}
		comment.Req = r
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return comments
}
