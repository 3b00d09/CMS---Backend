package handlers

import (
	"CMS-Backend/auth"
	"CMS-Backend/database"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)


type LoginData struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	
	var RegisterData database.UserCredentials

	err := c.BodyParser(&RegisterData)

	if err != nil{
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"message": "Internal Server Error.",
		})
	}

	validate := validator.New()

	err = validate.Struct(RegisterData)

	if err != nil{
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"message": "Incomplete form submission",
		})
	}

	isUniqueUsername, err := auth.IsUniqueUsername(RegisterData.Username)

	if(err != nil || !isUniqueUsername){
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userID, err := auth.CreateUser(RegisterData)

	if(err != nil){
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	cookie, err := auth.CreateSession(userID)

	if(err != nil){
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	c.Cookie(cookie)

	return c.JSON(fiber.Map{
		"username": RegisterData.Username,
	})
}

func Login(c *fiber.Ctx) error {

	var LoginData LoginData

	err := c.BodyParser(&LoginData)

	if err != nil{
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"message": "Internal Server Error.",
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
			"message": err.Error(),
		})
	}

	if(validUser.ID == ""){
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"message": "User does not exist",
		})
	}

	cookie, err := auth.CreateSession(validUser.ID)

	if(err != nil){
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"message": err.Error(),
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
			"message": err.Error(),
		})
	}
	
	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Now().Add(-time.Hour),
		Secure:   true,
		SameSite: "none",
	})

	return c.JSON(fiber.Map{
			"message": "Logout Successful.",
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