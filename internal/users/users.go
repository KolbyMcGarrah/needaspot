package users

import (
	"database/sql"
	"log"
	"strconv"

	database "github.com/KolbyMcGarrah/nas/internal/pkg/db/postgres"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	Level    string `json:"level"`
}

func (user *User) Create() {
	statement, err := database.Db.Prepare("INSERT INTO users(username,password,age,gender,level) VALUES($1,$2,$3,$4,$5)")
	log.Println(statement)
	if err != nil {
		log.Fatal("Error preparing create user sql: ", err)
	}
	hashedPassword, err := HashPassword(user.Password)
	_, err = statement.Exec(user.Username, hashedPassword, user.Age, user.Gender, user.Level)
	if err != nil {
		log.Fatal("Error executing create user SQL: ", err)
	}
}

//HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//Checkpassword hash compares raw password with it's hashed values
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//GetIDByUsername checks if a user exits in the database by a given username
func GetIDByUsername(username string) (int, error) {
	//Prepare the query
	statement, err := database.Db.Prepare("select user_id from users where username = $1")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(username)

	var Id int
	err = row.Scan(&Id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}

	return Id, nil
}

//GetUserByID returns user information from a given id
func GetUserById(id int) User {
	//prepare statement
	stmt, err := database.Db.Prepare("SELECT username, gender, level, age FROM users WHERE user_id=$1")
	if err != nil {
		log.Fatalf("Error grabbing user from id: %s", err)
	}
	row := stmt.QueryRow(id)
	var user User
	err = row.Scan(&user.Username, &user.Gender, &user.Level, &user.Age)
	if err != nil {
		log.Panicf("Failed to retrieve user with id: %v due to %s", id, err)
	}
	user.ID = strconv.Itoa(id)
	return user
}

func (user *User) Authenticate() bool {
	statement, err := database.Db.Prepare("select password from users WHERE Username = $1")
	if err != nil {
		log.Fatalf("Error preparing Authenticate sql: %s", err)
	}
	row := statement.QueryRow(user.Username)

	var hashedPassword string
	err = row.Scan(&hashedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}
	return CheckPassword(user.Password, hashedPassword)
}
