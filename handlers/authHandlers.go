package handlers

import (
	"CMS-Backend/auth"
	"CMS-Backend/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {

	
	user := database.UserCredentials{
		Username: "test",
		Password: "test",
	}

	passwordRepeat := "test"

	if(user.Password != passwordRepeat){
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"error": "Passwords do not match",
		})
	}

	isUniqueUsername, err := auth.IsUniqueUsername(user.Username)

	if(err != nil || !isUniqueUsername){
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userID, err := auth.CreateUser(user)

	if(err != nil){
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	cookie, err := auth.CreateSession(userID)

	if(err != nil){
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(cookie)

	return c.JSON(fiber.Map{
			"message": "Hello, Register!",
	})
}

func Login(c *fiber.Ctx) error {
	user := database.UserCredentials{
		Username: "test",
		Password: "test",
	}

	validUser, err := auth.UserExists(user)

	if(err != nil){
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if(validUser.ID == ""){
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"error": "User does not exist",
		})
	}

	cookie, err := auth.CreateSession(validUser.ID)

	if(err != nil){
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(cookie)

	return c.JSON(fiber.Map{
			"message": "Hello, Login!",
	})
}

func Logout(c *fiber.Ctx) error {

	cookie := c.Cookies("session_token")
	if(cookie == ""){
		c.Status(fiber.StatusUnprocessableEntity)
		return c.Redirect("/")
	}

	err := auth.ClearSession(cookie)

	if(err != nil){
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Now().Add(-time.Hour),
		Secure:   true,
		SameSite: "lax",
	})

	return c.JSON(fiber.Map{
			"message": "Hello, Logout!",
	})
}
