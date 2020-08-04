package users

import (
	"database/sql"
	"fmt"
	"testing"

	database "github.com/KolbyMcGarrah/nas/internal/pkg/db/postgres"
)

//TestCreate tests that we can successfully add users to the database
func TestCreate(t *testing.T) {
	username, check := InitializeData()
	if check != "" {
		t.Errorf(check)
	}
	check = CleanUp(username)
	if check != "" {
		t.Errorf(check)
	}
}

//TestPasswordHash tests both hashing and comparing hashed passwords.
func TestPasswordHash(t *testing.T) {
	tp := "testPassword"
	htp, err := HashPassword(tp)
	if err != nil {
		t.Errorf("Hash function did not execute correctly: %s", err)
	}
	if ok := CheckPassword("testPassword", htp); !ok {
		t.Error("Password hashes do not match.")
	}
}

//TestGetUserByID tests that we can get user data from the database by querying with the id
func TestGetUserByID(t *testing.T) {
	username, check := InitializeData()
	if check != "" {
		t.Fatalf("Error initializing test data: %s", check)
	}
	user_id, err := GetIDByUsername(username)
	if err != nil {
		t.Fatalf("Error retrieving test data: %s", err)
	}
	user := GetUserById(user_id)
	if username != user.Username {
		t.Fatalf("Usernames do not match. Have %s, expecting %s", user.Username, username)
	}
	check = CleanUp(username)
	if check != "" {
		t.Fatalf(check)
	}
}

//TestAuthenticateTrue tests that we can authenticate users with a username and password.
func TestAuthenticateTrue(t *testing.T) {
	username, check := InitializeData()
	if check != "" {
		t.Fatalf(check)
	}
	user_id, err := GetIDByUsername(username)
	if err != nil {
		t.Fatalf("Error retrieving test data: %s", err)
	}
	user := GetUserById(user_id)
	user.Password = "123Password"
	authenticated := user.Authenticate()

	if !authenticated {
		t.Fatalf("Failed to authenticate user.")
	}

	check = CleanUp(username)
	if check != "" {
		t.Fatalf(check)
	}
}

//TestAuthenticateFalse verifies that we do not authenticate users that do not exist in the database.
func TestAuthenticateFalse(t *testing.T) {
	username, check := InitializeData()
	if check != "" {
		t.Fatalf(check)
	}
	user_id, err := GetIDByUsername(username)
	if err != nil {
		t.Fatalf("Error retrieving user: %s", err)
	}
	user := GetUserById(user_id)
	user.Password = "alsdjfak"
	authenticate := user.Authenticate()
	if authenticate {
		t.Fatalf("Authenticated unregistered user.")
	}
	check = CleanUp(username)
	if check != "" {
		t.Fatalf(check)
	}
}

//InitializeData adds all necesary test cases to the database.
func InitializeData() (string, string) {
	//need to initialize database
	database.InitDB()
	//create a new user struct
	var user User
	user.Username = "testUser"
	var check string
	//check to see if username already exists, fail test if it does
	id, err := GetIDByUsername(user.Username)
	if id != 0 {
		check = fmt.Sprintf("Test username not removed from previous test.")
		return "", check
	} else if err != nil && err != sql.ErrNoRows {
		check = fmt.Sprintf("Received error when checking username %s", err)
		return "", check
	}
	//set values for all of the parameters to ensure that they can be added without an error.
	user.Age = 18
	user.Gender = "F"
	user.Level = "Intermediate"
	user.Password = "123Password"
	user.Create()
	return user.Username, check
}

//CleanUp removes test data from the database after a test completes.
func CleanUp(uname string) string {
	var check string
	stmt, err := database.Db.Prepare("DELETE FROM users WHERE user_id = $1")
	if err != nil {
		check = fmt.Sprintf("Error preparing cleanup sql stmt: %s", err)
	}
	tid, err := GetIDByUsername(uname)
	if err != nil {
		check = fmt.Sprintf("Error getting test username. Will need to manually delete this user.")
	}
	_, err = stmt.Exec(tid)
	if err != nil {
		check = fmt.Sprintf("Error cleaning up test case.")
	}
	return check
}
