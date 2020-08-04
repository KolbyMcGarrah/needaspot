package comments

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	database "github.com/KolbyMcGarrah/nas/internal/pkg/db/postgres"
	"github.com/KolbyMcGarrah/nas/internal/requests"
	"github.com/KolbyMcGarrah/nas/internal/users"
)

func TestCreate(t *testing.T) {
	comment, check := InitializeData()
	if check != "" {
		t.Fatalf(check)
	}
	check = comment.CleanUp()
	if check != "" {
		t.Fatalf(check)
	}
}

func TestRetrieve(t *testing.T) {
	comment, check := InitializeData()
	if check != "" {
		t.Fatalf(check)
	}
	req := comment.Req
	comments := RetrieveRequestComments(req)
	if len(comments) != 1 {
		t.Errorf("Invalid number of comments returned for request. Expected 1 but got %v", len(comments))
	}
	check = comment.CleanUp()
	if check != "" {
		t.Fatalf(check)
	}
}

//InitializeData prepares all test data for the unit tests.
func InitializeData() (Comment, string) {
	database.InitDB()
	// need to create a test user for the database
	var check string
	var user users.User
	user.Username = "testUser"
	//check to see if username already exists, fail test if it does
	id, err := users.GetIDByUsername(user.Username)
	if id != 0 {
		check = fmt.Sprintf("Test username not removed from previous test.")
		return Comment{}, check
	} else if err != nil && err != sql.ErrNoRows {
		check = fmt.Sprintf("Received error when checking username %s", err)
		return Comment{}, check
	}
	//set dummy values for the user
	user.Age = 18
	user.Gender = "F"
	user.Level = "Intermediate"
	user.Password = "123Password"
	user.Create()
	tempID, err := users.GetIDByUsername(user.Username)
	user.ID = strconv.Itoa(tempID)

	//Use the above user for creating the test Request for the comment
	var request requests.Request
	request.Location = "Yellowstone Fitness"
	request.Title = "Big Chest Day"
	request.Workout = "Chest"
	request.User = &user
	req_id := request.Save()
	request.ID = strconv.Itoa(req_id)

	//now create a comment for the request and use the user as the one making the comment.
	var comment Comment

	comment.Req = request
	comment.User = user
	comment.Text = "Test comment"
	cmt_id := comment.Save()
	comment.ID = strconv.Itoa(cmt_id)
	return comment, check
}

//CleanUp removes all test data from the database.
func (cmt *Comment) CleanUp() string {
	var check string
	//Clean up test cases
	//Clean up comment
	cmt_clean, err := database.Db.Prepare("DELETE FROM comments WHERE comment_id=$1")
	cmt_id, _ := strconv.Atoi(cmt.ID)
	usr_id, _ := strconv.Atoi(cmt.User.ID)
	req_id, _ := strconv.Atoi(cmt.Req.ID)
	if err != nil {
		check = fmt.Sprintf("Error when preparing comment cleanup: %s", err)
		return check
	}
	_, err = cmt_clean.Exec(cmt_id)
	if err != nil {
		check = fmt.Sprintf("Error executing comment cleanup sql. Manual clean-up required for ID %v: %s", cmt_id, err)
		return check
	}
	//Clean up request
	req_clean, err := database.Db.Prepare("DELETE FROM requests WHERE request_id=$1")
	if err != nil {
		check = fmt.Sprintf("Error preparing request cleanup sql: %s", err)
		return check
	}
	_, err = req_clean.Exec(req_id)
	if err != nil {
		check = fmt.Sprintf("Error cleaning up test request with id: %v: %s", req_id, err)
		return check
	}
	//Clean up user
	usr_clean, err := database.Db.Prepare("DELETE FROM users WHERE user_id=$1")
	if err != nil {
		check = fmt.Sprintf("Error preparing cleanup sql for test user: %s", err)
		return check
	}
	_, err = usr_clean.Exec(usr_id)
	if err != nil {
		check = fmt.Sprintf("Error deleting user with ID: %v. Manual cleanup required. Error: %s", usr_id, err)
		return check
	}
	return check
}
