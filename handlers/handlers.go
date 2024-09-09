package handlers

import (
	"CMS-Backend/auth"
	"CMS-Backend/database"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type CreateProjectData struct {
	ProjectName string `json:"project_name" validate:"required"`
	ProjectDescription string `json:"project_description"`
}

func HandleCreateProject(c fiber.Ctx) error {
	cookie := c.Cookies("session_token")
	user := auth.AuthenticateSession(cookie)

	if(len(user.Username) == 0){
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Invalid session.",
		})
	}

	var ProjectData CreateProjectData

	c.Bind().Body(&ProjectData)

	validate := validator.New()

	err := validate.Struct(ProjectData)

	if err != nil{
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"message": "Incomplete form submission",
		})
	}

	statement, err := database.DB.Prepare("INSERT INTO projects (creator_id, name, description) VALUES (?, ?, ?)")

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error.",
			"detail":err.Error(),
		})
	}

	defer statement.Close()

	_, err = statement.Exec(user.ID, ProjectData.ProjectName, ProjectData.ProjectDescription)

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error.",
			"detail": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Project created successfully.",
		"data": fiber.Map{
			"name":        ProjectData.ProjectName,
			"description": ProjectData.ProjectDescription,
		},
	})

}

func HandleGetProjects(c fiber.Ctx) error{
	cookie := c.Cookies("session_token")
	user := auth.AuthenticateSession(cookie)

	if(len(user.Username) ==  0 || len(user.ID) == 0){
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Invalid session.",
		})
	}

	statement, err := database.DB.Prepare("SELECT name, description FROM projects WHERE creator_id = ?")

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
			"detail": err.Error(),
		})
	}

	defer statement.Close()

	var ProjectData []struct {
		ProjectName        string
		ProjectDescription string
	}

	rows, err := statement.Query(user.ID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
			"detail":  err.Error(),
		})
	}
	defer rows.Close()

	for rows.Next() {
		var project struct {
			ProjectName        string
			ProjectDescription string
		}
		if err := rows.Scan(&project.ProjectName, &project.ProjectDescription); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "Internal server error",
				"detail":  err.Error(),
			})
		}
		ProjectData = append(ProjectData, project)
	}

	return c.JSON(fiber.Map{
		"projects": ProjectData,
	})
}