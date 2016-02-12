package main

import (
	// Standard library packages
	"net/http"

	// Third party packages
	"github.com/gin-gonic/gin"
	"./controllers"
	"gopkg.in/mgo.v2"
)

func main() {

	// Get a UserController instance
	uc := controllers.NewUserController(getSession())

	// Get a user resource
	router := gin.Default()
    router.GET("/users", uc.usersList)
    router.GET("/users/:id", uc.usersDetail)
    router.GET("/usersdelete/:id", uc.usersDetail)

    router.POST("/usersadd", ucuserPost)
    router.POST("/usersupdate", ucuserPost)
    
    router.Run(":8000")
}

// getSession creates a new mongo session and panics if connection error occurs
func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}

	// Deliver session
	return s
}
