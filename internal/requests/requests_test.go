package requests

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	database "github.com/KolbyMcGarrah/nas/internal/pkg/db/postgres"
	"github.com/KolbyMcGarrah/nas/internal/users"
)

func TestSave(t *testing.T) {
	user, request, check := InitializeData()
	if check != "" {
		t.Fatalf(check)
	}
	req_id := request.Save()

	check = CleanUp(user.Username, req_id)
	if check != "" {
		t.Fatalf(check)
	}
}

func TestGetReqByID(t *testing.T) {
	user, request, check := InitializeData()
	if check != "" {
		t.Fatalf(check)
	}
	req_id := request.Save()

	// Test Get Reuqest by ID

	found := GetReqByID(req_id)
	if found.Location != request.Location {
		t.Errorf("Expected %s but recieved %s", request.Location, found.Location)
	}
	check = CleanUp(user.Username, req_id)
	if check != "" {
		t.Fatalf(check)
	}
}

//TestGetAll ensures that the GetAll function returns a slice of requests for each request in the database. (We only check that our added test case exists, however)
func TestGetAll(t *testing.T) {
	user, request, check := InitializeData()
	if check != "" {
		t.Fatalf(check)
	}
	req_id := request.Save()

	//Get all of the requests and make sure our added test case is in there
	var found bool
	rSlice := GetAll()
	for _, req := range rSlice {
		if req.ID == strconv.Itoa(req_id) {
			found = true
		}
		if !found {
			t.Errorf("Error, did not retrieve all requests. Missing added test case.")
		}
	}

	check = CleanUp(user.Username, req_id)
	if check != "" {
		t.Fatalf(check)
	}
}

//InitializeData sets up the test data needed to run the unit tests
func InitializeData() (users.User, Request, string) {
	//Initialize database
	database.InitDB()

	// need to create a test user for the database
	var user users.User
	user.Username = "testUser"
	//check to see if username already exists, fail test if it does
	id, err := users.GetIDByUsername(user.Username)
	if id != 0 {
		return users.User{}, Request{}, "Test username not removed from previous test."
	} else if err != nil && err != sql.ErrNoRows {
		es := fmt.Sprintf("Received error when checking username %s", err)
		return users.User{}, Request{}, es
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
	var request Request
	request.Location = "Yellowstone Fitness"
	request.Title = "Big Chest Day"
	request.Workout = "Chest"
	request.User = &user
	request.Time = "2020-10-12 11:30"

	return user, request, ""
}

//CleanUp removes the test data from the database so that tests are re-runable.
func CleanUp(username string, req_id int) string {
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
	return ""
}
