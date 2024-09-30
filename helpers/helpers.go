package helpers

import (
	"CMS-Backend/database"
	"fmt"
	"strings"
	"time"
)

func UnixToHuman(unix int64) string {
	currTime := time.Now().Unix()

	// get the difference in seconds
	diff := currTime - unix
	if diff < 60 {
		return "Just now"
	} else if diff < 3600 {
		return "Less than an hour ago"
	} else if diff < 86400 {
		return fmt.Sprintf("%d hours ago", diff / 3600)
	}else{
		return fmt.Sprintf("%d days ago", diff / 86400)
	}

}


func GetProjectIdByName(name string, userID string) string{
	
	// if passing in invalid query param
	if name == ""{
		return ""
	}
	statement, err := database.DB.Prepare("SELECT id FROM projects WHERE LOWER(name) = ? AND creator_id = ?")

	if err != nil{
		return ""
	}

	defer statement.Close()

	var projectId string

	err = statement.QueryRow(strings.ToLower(name), userID).Scan(&projectId)

	if err != nil{
		return ""
	} 

	return projectId
	
}