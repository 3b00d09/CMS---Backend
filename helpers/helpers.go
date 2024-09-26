package helpers

import (
	"fmt"
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

