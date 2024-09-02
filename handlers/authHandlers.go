package handlers

import (
	"CMS-Backend/auth"
	"CMS-Backend/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RegisterData struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeatPassword"`
}

type LoginData struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	
	var RegisterData RegisterData

	err := c.BodyParser(&RegisterData)

	if err != nil{
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"error": "Internal Server Error.",
		})
	}


	if(RegisterData.Password != RegisterData.RepeatPassword){
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"error": "Passwords do not match",
		})
	}

	isUniqueUsername, err := auth.IsUniqueUsername(RegisterData.Username)

	if(err != nil || !isUniqueUsername){
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user := database.UserCredentials{
		Username: RegisterData.Username,
		Password: RegisterData.Password,
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
		"username": user.Username,
	})
}

func Login(c *fiber.Ctx) error {

	var LoginData LoginData

	err := c.BodyParser(&LoginData)

	if err != nil{
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"error": "Internal Server Error.",
		})
	}

	user := database.UserCredentials{
		Username: LoginData.Username,
		Password: LoginData.Password,
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
		"username": user.Username,
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
		Secure:   false,
		SameSite: "lax",
	})

	return c.JSON(fiber.Map{
			"message": "Hello, Logout!",
	})
}


func ValidateSession(c *fiber.Ctx)(error){
	cookie := c.Cookies("session_token")
	if(cookie == ""){
		return c.JSON(fiber.Map{
				"username": "",
		})
	}
	user := auth.AuthenticateSession(cookie)
		
	return c.JSON(fiber.Map{
			"username": user.Username,
	})
}