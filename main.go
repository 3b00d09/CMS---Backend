package main

import (
	"fmt"
)

const port string = ":3000"

func main() {
	app := SetupRoutes()
	app.Listen(port)
	fmt.Printf("Server Running on http://localhost%s\n", port)
}