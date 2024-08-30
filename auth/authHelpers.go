package auth

import (
	"CMS-Backend/database"
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)


func IsUniqueUsername(username string) (bool, error) {
	statement, err := database.DB.Prepare("SELECT username FROM user WHERE username = ?")
	if err != nil {
		return false, fmt.Errorf("internal server error")
	}
	defer statement.Close()
	var user string

	err = statement.QueryRow(username).Scan(&user)
	if err != nil {
		// no rows means the user doesnt exist so the username is unique
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, fmt.Errorf("internal server error")
	}
	// if we reach this point, the user exists and thus not unique
	return false, fmt.Errorf("username already taken")
}


func GeneratHashedPassword(password string) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(err)
	}

	return hashedPassword
}

func CheckPasswordHash(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))

	return err == nil
}