package handlers

import (
	"CMS-Backend/auth"
	"CMS-Backend/database"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CreateProjectData struct {
	ProjectName string `json:"project_name" db:"name" validate:"required"`
}

func HandleCreateProject(c *fiber.Ctx) error {
	cookie := c.Cookies("session_token")
	user := auth.AuthenticateSession(cookie)

	if(len(user.Username) == 0){
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Invalid session.",
		})
	}

	var ProjectData CreateProjectData

	c.BodyParser(&ProjectData)

	validate := validator.New()

	err := validate.Struct(ProjectData)

	if err != nil{
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"message": "Incomplete form submission",
		})
	}

	statement, err := database.DB.Prepare("INSERT INTO projects (creator_id, name) VALUES (?, ?)")

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error.",
		})
	}

	defer statement.Close()

	_, err = statement.Exec(user.ID, ProjectData.ProjectName)

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error.",
			"detail": err.Error(),
		})
	}

	return c.JSON((fiber.Map{
		"message": "Project created successfully.",
	}))

}