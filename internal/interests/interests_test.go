package interests

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	database "github.com/KolbyMcGarrah/nas/internal/pkg/db/postgres"
	"github.com/KolbyMcGarrah/nas/internal/requests"

	"github.com/KolbyMcGarrah/nas/internal/users"
)

//TestSave tests that we can save an interest to the database.
func TestSave(t *testing.T) {
	//Initialize the database with a user, and request. Check for errors with the check var.
	u, r, c := InitializeData()
	if c != "" {
		t.Fatalf(c)
	}
	var i Interest
	i.User = &u
	i.Request = &r
	i.Description = "I would like to work out with you"
	i.ID = i.Save()
	req_id, _ := strconv.Atoi(i.Request.ID)
	c = CleanUp(i.User.Username, req_id, i.ID)
}

//TestGetInterestsByReqID tests that we can retrieve all interest from a given request.
func TestGetInterestsByReqID(t *testing.T) {
	//Initialize the database with a user, and request. Check for errors with the check var.
	u, r, c := InitializeData()
	if c != "" {
		t.Fatalf(c)
	}
	var i Interest
	i.User = &u
	i.Request = &r
	i.Description = "I would like to work out with you"
	i.ID = i.Save()
	interests := GetInterestByReqID(r)
	if len(interests) < 1 {
		t.Fatalf("No interests returned. %v", len(interests))
	}
	for _, interest := range interests {
		if interest.User.Username != u.Username {
			t.Fatalf("Username mismatch. Have %s, but expecting %s", i.User.Username, u.Username)
		}
	}
	req_id, _ := strconv.Atoi(r.ID)
	//Cleanup test data.
	c = CleanUp(u.Username, req_id, i.ID)
	if c != "" {
		t.Fatal(c)
	}
}

//TestGetInterestByUser tests that we can query all interests made by a user.
func TestGetINterestByUser(t *testing.T) {
	//Initialize the database with a user, and request. Check for errors with the check var.
	u, r, c := InitializeData()
	if c != "" {
		t.Fatalf(c)
	}
	//Create and save an interest in the returned request
	var i Interest
	i.User = &u
	i.Request = &r
	i.Description = "I would like to work out with you"
	i.ID = i.Save()

	interests := GetInterestsByUser(u)

	if len(interests) < 1 {
		t.Fatalf("No interests returned. %v", len(interests))
	}

	for _, interest := range interests {
		if interest.User.Username != u.Username {
			t.Fatalf("Username mismatch. Have %s, but expecting %s", i.User.Username, u.Username)
		}
	}
	//Cleanup test data
	req_id, _ := strconv.Atoi(r.ID)
	c = CleanUp(u.Username, req_id, i.ID)
	if c != "" {
		t.Fatal(c)
	}
}

//TestAccept ensures that we can update an interest record to accepted
func TestAccept(t *testing.T) {
	//Initialize the database with a user, and request. Check for errors with the check var.
	u, r, c := InitializeData()
	if c != "" {
		t.Fatalf(c)
	}
	//Create and save an interest in the returned request
	var i Interest
	i.User = &u
	i.Request = &r
	i.Description = "I would like to work out with you"
	i.ID = i.Save()
	i.Accept(u)
	if i.AcceptedTs == "" {
		t.Error("Error, accepted timestamp not updated correctly.")
	}
	if i.AcceptedUser == nil {
		t.Error("Error, accepted user is not updated.")
	}
	if !i.Accepted {
		t.Error("Error, accepted boolean not updated to true")
	}
	req_id, _ := strconv.Atoi(r.ID)
	c = CleanUp(u.Username, req_id, i.ID)
	if c != "" {
		t.Fatalf(c)
	}
}

//InitializeData sets up the test data needed to run the unit tests
func InitializeData() (users.User, requests.Request, string) {
	//Initialize database
	database.InitDB()

	// need to create a test user for the database
	var user users.User
	user.Username = "testUser"
	//check to see if username already exists, fail test if it does
	id, err := users.GetIDByUsername(user.Username)
	if id != 0 {
		return users.User{}, requests.Request{}, "Test username not removed from previous test."
	} else if err != nil && err != sql.ErrNoRows {
		es := fmt.Sprintf("Received error when checking username %s", err)
		return users.User{}, requests.Request{}, es
	}
	//set dummy values for the user
	user.Age = 18
	user.Gender = "F"
	user.Level = "Intermediate"
	user.Password = "123Password"
	user.Create()
	tempID, err := users.GetIDByUsername(user.Username)
	user.ID = strconv.Itoa(tempID)

	//Use the above user for creating the Request
	var request requests.Request
	request.Location = "Yellowstone Fitness"
	request.Title = "Big Chest Day"
	request.Workout = "Chest"
	request.User = &user
	req_id := request.Save()
	request.ID = strconv.Itoa(req_id)

	return user, request, ""
}

func CleanUp(username string, req_id int, interest_id int) string {
	var check string

	//Cleanup test.
	//delete entry from request table
	stmt, err := database.Db.Prepare("delete from requests where request_id=$1")
	if err != nil {
		check = fmt.Sprintf("Error preparing cleanup sql. Will need to manually clear row :%v", req_id)
		return check
	}
	_, err = stmt.Exec(req_id)
	if err != nil {
		check = fmt.Sprintf("Error cleaning up test. Please manually delete row with id: %v from request table.", req_id)
		return check
	}

	//delete entry from user table
	stmt, err = database.Db.Prepare("DELETE FROM users WHERE user_id = $1")
	if err != nil {
		check = fmt.Sprintf("Error preparing cleanup sql stmt: %s", err)
		return check
	}
	tid, err := users.GetIDByUsername(username)
	if err != nil {
		check = ("Error getting test username. Will need to manually delete this user.")
		return check
	}
	_, err = stmt.Exec(tid)
	if err != nil {
		check = ("Error cleaning up test case.")
		return check
	}

	//delete entry from interest Table
	stmt, err = database.Db.Prepare("DELETE FROM interests WHERE interest_id=$1")
	if err != nil {
		check = fmt.Sprintf("Error preparing sql to delete interest from database with ID: %v", interest_id)
	}
	_, err = stmt.Exec(interest_id)
	if err != nil {
		check = fmt.Sprintf("Error deleting interest from database with ID: %v", interest_id)
	}
	return ""
}
