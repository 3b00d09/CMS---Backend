package auth

import (
	"CMS-Backend/database"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func AuthenticateSession(cookie string) database.User {

	statement, err := database.DB.Prepare("SELECT id, user_id, active_expires FROM user_session WHERE id = ?")

	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	row := statement.QueryRow(cookie)

	var sessionID, userID string
	var activeExpires int64

	err = row.Scan(&sessionID, &userID, &activeExpires)
	if err != nil {
		return database.User{}
	}

	 if(activeExpires < time.Now().Unix()){
		return database.User{}
	 }

	 statement, err = database.DB.Prepare("SELECT username FROM user WHERE id = ?")

	 if err != nil {	
		 log.Fatal(err)
	 }
	 defer statement.Close()

	 row = statement.QueryRow(userID)
	 User := database.User{}
	 err = row.Scan(&User.Username)

	 if err != nil {
		 return database.User{}
	 }

	 User.ID = userID
	 return User	

}

func UserExists(User database.UserCredentials) (database.User, error) {

	statement, err := database.DB.Prepare("SELECT id, username FROM user WHERE username = ?")
	if err != nil {
		return database.User{}, fmt.Errorf("internal server error")
	}

	defer statement.Close()


	var user database.User
	err = statement.QueryRow(User.Username).Scan(&user.ID, &user.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User doesnt exist")
			return  database.User{}, nil
		}
		log.Fatal(err)
	}

	statement, err = database.DB.Prepare("SELECT password FROM user WHERE username = ?")
	if err != nil {
		log.Fatal(err)
	}

	var password []byte
	err = statement.QueryRow(User.Username).Scan(&password)
	if err != nil {
		log.Fatal(err)
	}

	if !CheckPasswordHash(User.Password, []byte(password)) {
		return  database.User{}, nil
	}

	return user, nil

}

func CreateUser(user database.UserCredentials) (string, error) {
	hashedPassword := GeneratHashedPassword(user.Password)
	
	statement, err := database.DB.Prepare("INSERT INTO user (id, username, password, first_name, last_name, email) VALUES (?, ?, ?, ?, ?, ?)")
	
	if err != nil {
		return "", fmt.Errorf("internal server error")
	}

	defer statement.Close()

	user.ID = uuid.New().String()

	_, err = statement.Exec(user.ID, user.Username, hashedPassword, user.FirstName, user.LastName, user.Email)

	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("internal server error")
	}

	return user.ID, nil

}



func CreateSession(userId string) (*fiber.Cookie, error) {

	sessionId := uuid.New().String()

	newSession := database.UserSession{
		ID:            sessionId,
		UserID:        userId,
		ActiveExpires: time.Now().Add(3600 * time.Hour * 24 * 7).Unix(),
		IdleExpires:   0,
	}
	
	statement, err := database.DB.Prepare("INSERT INTO user_session (id, user_id, active_expires, idle_expires) VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("session creation failed")
	}

	defer statement.Close()

	_, err = statement.Exec(newSession.ID, newSession.UserID, newSession.ActiveExpires, newSession.IdleExpires)
	if err != nil {
		return nil, fmt.Errorf("session creation failed")
	}


	cookie := &fiber.Cookie{
		Name:     "session_token",
		Value:    sessionId,
		Path:     "/",
		MaxAge:   int(newSession.ActiveExpires),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "none",
		Partitioned: true,
		
	}

	return cookie, nil

}

func ClearSession(token string) (error){
	statement, err := database.DB.Prepare("DELETE FROM user_session WHERE id = ?")
	
	if err != nil{
		return fmt.Errorf("internal server error")
	}
	defer statement.Close()

	_, err = statement.Exec(token)

	if err != nil{
		return fmt.Errorf("internal server error")
	}

	return nil

}

func ClearUserSessions(userId string){
	statement, err := database.DB.Prepare("DELETE FROM user_session WHERE user_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	_, err = statement.Exec(userId)

	if err != nil {
		fmt.Print(err)
	}
}