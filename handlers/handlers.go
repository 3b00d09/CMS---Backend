package handlers

import (
	"CMS-Backend/auth"
	"CMS-Backend/database"
	"CMS-Backend/helpers"
	"strings"

	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type CreateProjectData struct {
	ProjectName string `json:"project_name" validate:"required"`
	ProjectDescription string `json:"project_description"`
}

type CreatePageData  struct {
	ProjectName string `json:"project_name" validate:"required"`
	PageName string `json:"page_name" validate:"required"`
}

type GetPageData struct {
	ProjectName string `json:"project_name" validate:"required"`
	PageName string `json:"page_name" validate:"required"`
}

type UpdatePageData struct {
	ProjectName string `json:"project_name" validate:"required"`
	PageName string `json:"page_name" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func HandleCreateProject(c fiber.Ctx) error {
	cookie := c.Cookies("session_token")
	user := auth.AuthenticateSession(cookie)

	if(len(user.Username) == 0){
		c.Status(fiber.StatusUnauthorized)
		return helpers.SessionError(c)
	}

	var ProjectData CreateProjectData

	c.Bind().Body(&ProjectData)

	validate := validator.New()

	err := validate.Struct(ProjectData)

	if err != nil{
		c.Status(fiber.StatusUnprocessableEntity)
		return helpers.FormError(c)
	}

	statement, err := database.DB.Prepare("INSERT INTO projects (creator_id, name, description, last_updated) VALUES (?, ?, ?, ?)")

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	defer statement.Close()

	_, err = statement.Exec(user.ID, ProjectData.ProjectName, ProjectData.ProjectDescription, time.Now().Unix())

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
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
		return helpers.SessionError(c)
	}

	includePages := c.Query("pages")

	if includePages == "true"{
		return HandleGetProjectWithPages(c)
	}

	statement, err := database.DB.Prepare("SELECT name, description FROM projects WHERE creator_id = ?")

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	defer statement.Close()

	var ProjectData []struct {
		ProjectName        string
		ProjectDescription string
	}

	rows, err := statement.Query(user.ID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}
	defer rows.Close()

	for rows.Next() {
		var project struct {
			ProjectName        string
			ProjectDescription string
		}
		if err := rows.Scan(&project.ProjectName, &project.ProjectDescription); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return helpers.ServerError(c, err)
		}
		ProjectData = append(ProjectData, project)
	}

	return c.JSON(fiber.Map{
		"projects": ProjectData,
	})
}

func HandleSearchUsers(c fiber.Ctx) error{
	cookie := c.Cookies("session_token")
	user := auth.AuthenticateSession(cookie)

	if(len(user.Username) ==  0 || len(user.ID) == 0){
		c.Status(fiber.StatusUnauthorized)
		return helpers.SessionError(c)
	}

	queryParam := c.Query("q")
	queryParam = strings.TrimSpace(queryParam)

	if(len(queryParam) == 0){
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"data": []string{},
		})
	}

	statement, err := database.DB.Prepare("SELECT username FROM user WHERE LOWER(username) LIKE LOWER(?) AND id != ?")

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	defer statement.Close()

	var users []string

	rows, err := statement.Query("%" + queryParam + "%", user.ID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return helpers.ServerError(c, err)
		}
		users = append(users, username)
	}

	return c.JSON(fiber.Map{
		"data": users,
	})

}

func HandleCreatePage(c fiber.Ctx) error{
	cookie := c.Cookies("session_token")
	user := auth.AuthenticateSession(cookie)

	if(len(user.Username) ==  0 || len(user.ID) == 0){
		c.Status(fiber.StatusUnauthorized)
		return helpers.SessionError(c)
	}

	var PageData CreatePageData
	c.Bind().Body(&PageData)
	validator := validator.New()
	err := validator.Struct(PageData)

	if err != nil{
		c.Status(fiber.StatusUnprocessableEntity)
		return helpers.FormError(c)
	}

	statement, err := database.DB.Prepare("SELECT id, creator_id FROM projects WHERE name = ?")

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	defer statement.Close()
	
	var project struct{
		ProjectID string `db:"id"`
		CreatorID string `db:"creator_id"`
	}
	err = statement.QueryRow(PageData.ProjectName).Scan(&project.ProjectID, &project.CreatorID)

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	if(project.CreatorID != user.ID){
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "You are not authorized to create a page in this project.",
		})
	}

	statement, err = database.DB.Prepare("INSERT INTO pages (project_id, name) VALUES (?, ?)")

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}


	_, err = statement.Exec(project.ProjectID, PageData.PageName)

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	statement, err = database.DB.Prepare("UPDATE projects SET last_updated = ? WHERE id = ?")

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	_, err = statement.Exec(time.Now().Unix(), project.ProjectID)

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Page created successfully",
	})
}

func HandleGetPage(c fiber.Ctx) error{
	cookie := c.Cookies("session_token")
	user := auth.AuthenticateSession(cookie)
	if(len(user.Username) ==  0 || len(user.ID) == 0){
		c.Status(fiber.StatusUnauthorized)
		return helpers.SessionError(c)
	}

	pageQueryParam := c.Query("page")
	projectQueryParam := c.Query("project")

	if(len(pageQueryParam) == 0 || len(projectQueryParam) == 0){
		c.Status(fiber.StatusBadRequest)
		return helpers.FormError(c)
	}

	statement, err := database.DB.Prepare("SELECT id, creator_id FROM projects WHERE name = ?")

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	defer statement.Close()
	
	var project struct{
		ProjectID string `db:"id"`
		CreatorID string `db:"creator_id"`
	}

	err = statement.QueryRow(projectQueryParam).Scan(&project.ProjectID, &project.CreatorID)

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	if(project.CreatorID != user.ID){
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "You are not authorized to access this page.",
		})
	}

	statement, err = database.DB.Prepare("SELECT content FROM pages WHERE project_id = ? AND name = ?")

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	var content string
	err = statement.QueryRow(project.ProjectID, pageQueryParam).Scan(&content)

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	return c.JSON(fiber.Map{
		"content": content,
	})

}

func HandleUpdatePage(c fiber.Ctx) error {
	cookie := c.Cookies("session_token")
	user := auth.AuthenticateSession(cookie)

	if(len(user.Username) ==  0 || len(user.ID) == 0){
		c.Status(fiber.StatusUnauthorized)
		return helpers.SessionError(c)
	}

	var PageData UpdatePageData
	c.Bind().Body(&PageData)
	validator := validator.New()
	err := validator.Struct(PageData)

	if err != nil{
		c.Status(fiber.StatusUnprocessableEntity)
		return helpers.FormError(c)
	}

	statement, err := database.DB.Prepare("SELECT id, creator_id FROM projects WHERE name = ?")

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	defer statement.Close()
	
	var project struct{
		ProjectID string `db:"id"`
		CreatorID string `db:"creator_id"`
	}
	err = statement.QueryRow(PageData.ProjectName).Scan(&project.ProjectID, &project.CreatorID)

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	if(project.CreatorID != user.ID){
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "You are not authorized to edit a page in this project.",
		})
	}

	statement, err = database.DB.Prepare("UPDATE pages SET content = ? WHERE project_id = ? AND name = ?")

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	_, err = statement.Exec(PageData.Content, project.ProjectID, PageData.PageName)

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	statement, err = database.DB.Prepare("UPDATE projects SET last_updated = ? WHERE id = ?")

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	_, err = statement.Exec(time.Now().Unix(), project.ProjectID)

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Page updated successfully",
	})

}

func HandleGetProjectWithPages(c fiber.Ctx) error{
	cookie := c.Cookies("session_token")
	user := auth.AuthenticateSession(cookie)

	if(len(user.Username) ==  0 || len(user.ID) == 0){
		c.Status(fiber.StatusUnauthorized)
		return helpers.SessionError(c)
	}

	projectQueryParam := c.Query("project")

	if(len(projectQueryParam) == 0){
		c.Status(fiber.StatusBadRequest)
		return helpers.FormError(c)
	}
	

	statement, err := database.DB.Prepare("SELECT id FROM projects WHERE creator_id = ? AND name = ?")

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	defer statement.Close()

	var projectID string

	err = statement.QueryRow(user.ID, projectQueryParam).Scan(&projectID)

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}


	statement, err = database.DB.Prepare("SELECT pages.name FROM pages JOIN projects ON pages.project_id = projects.id WHERE projects.id = ?")

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	defer statement.Close()

	var pages []string

	rows, err := statement.Query(projectID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	defer rows.Close()
	
	for rows.Next() {
		var page struct{
			PageName string `db:"name"`
		}
		if err := rows.Scan(&page.PageName); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return helpers.ServerError(c, err)
		}
		pages = append(pages, page.PageName)
	}

	return c.JSON(fiber.Map{
		"pages": pages,
	})
}

func HandleGetStatsPage(c fiber.Ctx) error {
	cookie := c.Cookies("session_token")
	user := auth.AuthenticateSession(cookie)

	if(len(user.Username) ==  0 || len(user.ID) == 0){
		c.Status(fiber.StatusUnauthorized)
		return helpers.SessionError(c)
	}

	statement, err := database.DB.Prepare(`
    SELECT 
        (SELECT COUNT(id) FROM projects WHERE creator_id = ?) AS project_count,
        (SELECT COUNT(pages.id) 
         FROM pages 
         INNER JOIN projects ON pages.project_id = projects.id 
         WHERE projects.creator_id = ?) AS page_count
	`)

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	defer statement.Close()

	var counts struct{
		ProjectCount int `db:"project_count"`
		PageCount int `db:"page_count"`
	}

	err = statement.QueryRow(user.ID, user.ID).Scan(&counts.ProjectCount, &counts.PageCount)

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	return c.JSON(fiber.Map{
		"project_count": counts.ProjectCount,
		"page_count": counts.PageCount,
	})
}


func HandleGetLastModified(c fiber.Ctx) error{
	cookie := c.Cookies("session_token")
	user := auth.AuthenticateSession(cookie)

	if(len(user.Username) ==  0 || len(user.ID) == 0){
		c.Status(fiber.StatusUnauthorized)
		return helpers.SessionError(c)
	}

	statement, err := database.DB.Prepare("SELECT name, last_updated FROM projects WHERE creator_id = ? ORDER BY last_updated DESC LIMIT 3")

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	defer statement.Close()

	var projects []struct{
		ProjectName string `db:"name"`
		LastUpdated int `db:"last_updated"`
	}

	rows, err := statement.Query(user.ID)

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	defer rows.Close()

	for rows.Next() {
		var project struct {
			ProjectName string `db:"name"`
			LastUpdated int `db:"last_updated"`
		}
		err := rows.Scan(&project.ProjectName, &project.LastUpdated)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return helpers.ServerError(c, err)
		}
		projects = append(projects, project)
	}

	return c.JSON(fiber.Map{
		"projects": projects,
	})


}

func HandleDeletePage(c fiber.Ctx) error {
	user := auth.AuthenticateSession(c.Cookies("session_token"))

	if(user.Username == ""){
		c.Status(fiber.StatusUnauthorized)
		return helpers.SessionError(c)
	}

	pageQueryParam := c.Query("page")
	projectQueryParam := c.Query("project")

	if(len(pageQueryParam) == 0 || len(projectQueryParam) == 0){
		c.Status(fiber.StatusBadRequest)
		return helpers.FormError(c)
	}

	statement, err := database.DB.Prepare(`
		DELETE FROM pages
		WHERE pages.project_id = (SELECT id FROM projects WHERE creator_id = ? AND LOWER(name) = ?)
		AND LOWER(pages.name) = ?
	`)

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	defer statement.Close()

	result, err := statement.Exec(user.ID, strings.ToLower(projectQueryParam), strings.ToLower(pageQueryParam))

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	rowsAffected, err := result.RowsAffected()	

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	if rowsAffected == 0{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Page not found.",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Page deleted successfully",
	})
}


func HandleDeleteProject(c fiber.Ctx) error{
	user := auth.AuthenticateSession(c.Cookies("session_token"))

	if user.Username == ""{
		c.Status(fiber.StatusUnauthorized)
		return helpers.SessionError(c)
	}

	projectQueryParam := c.Query("project")

	if len(projectQueryParam) == 0{
		c.Status(fiber.StatusBadRequest)
		return helpers.FormError(c)
	}

	statement, err := database.DB.Prepare(`
		DELETE FROM projects
		WHERE creator_id = ? AND LOWER(name) = ?
	`)

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	defer statement.Close()

	result, err := statement.Exec(user.ID, strings.ToLower(projectQueryParam))

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	if rowsAffected == 0{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Project not found.",
		})
	}
	
	return c.JSON(fiber.Map{
		"message": "Project deleted successfully",
	})

}

func HandleCreateTodo(c fiber.Ctx) error{

	user := auth.AuthenticateSession(c.Cookies("session_token"))

	if(user.Username == ""){
		c.Status(fiber.StatusUnauthorized)
		return helpers.SessionError(c)
	}

	var todo struct{
		ProjectName string `json:"project_name" validate:"required"`
		Content string `json:"content" validate:"required"`
	}

	c.Bind().Body(&todo)
	validator := validator.New()
	err := validator.Struct(todo)


	if err != nil{
		c.Status(fiber.StatusUnprocessableEntity)
		return helpers.FormError(c)
	}


	projectId := helpers.GetProjectIdByName(todo.ProjectName, user.ID)

	if projectId == ""{
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Project not found",
		})
	}

	statement, err := database.DB.Prepare("INSERT INTO todo (project_id, content) VALUES (?, ?)")

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err);
	}

	defer statement.Close()

	
	_, err = statement.Exec(projectId, todo.Content)

	if err != nil{
		
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Todo successfully created",
	})

}

func HandleGetTodos(c fiber.Ctx) error{
	user := auth.AuthenticateSession(c.Cookies("session_token"))

	if user.Username == ""{
		c.Status(fiber.StatusUnauthorized)
		return helpers.SessionError(c)
	}

	if c.Query("project") == ""{
		c.Status(fiber.StatusUnprocessableEntity)
		return helpers.FormError(c)
	}
	projectID := helpers.GetProjectIdByName(c.Query("project"), user.ID)

	if(len(projectID) == 0){
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Project not found",
		})
	}

	statement, err := database.DB.Prepare("SELECT content, completed FROM todo WHERE project_id = ?")

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	defer statement.Close()


	var todos[] struct {
		Content string
		Completed bool
	}

	rows, err := statement.Query(projectID)

	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return helpers.ServerError(c, err)
	}

	for rows.Next(){
		var todo struct{
			Content string
			Completed bool
		}

		err := rows.Scan(&todo.Content, &todo.Completed)
		if err != nil{
			c.Status(fiber.StatusInternalServerError)
			return helpers.ServerError(c, err)
		}

		todos = append(todos, todo)
	}

	defer rows.Close()

	return c.JSON(fiber.Map{
		"data": todos,
	})

}