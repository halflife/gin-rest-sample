package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"../models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// UserController represents the controller for operating on the User resource
	UserController struct {
		session *mgo.Session
	}
)

const (
    DB_NAME       = "gotest"
    DB_COLLECTION = "pepole_new1"
)

// NewUserController provides a reference to a UserController with provided mongo session
func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

// Get all Users
func (uc UserController) usersList(c *gin.Context) {
	var results []models.User
    err = uc.session.DB(DB_NAME).C(DB_COLLECTION).Find(nil).All(&results)
	checkErr(err, "Users doesn't exist")
    content := gin.H{}
    for k, v := range users {
        content[strconv.Itoa(k)] = v
    }
    c.JSON(200, content)
}

// GetUser retrieves an individual user resource
func (uc UserController) GetUser(c *gin.Context) {
	// Grab id
	id := c.Params.ByName("id")
	
	 //Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		c.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Stub user
	u := models.User{}

	// Fetch user
	if err := uc.session.DB(DB_NAME).C(DB_COLLECTION).FindId(oid).One(&u); err != nil {
		checkErr(err, "Users doesn't exist")
		c.WriteHeader(403)

		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	c.Header().Set("Content-Type", "application/json")
    c.JSON(200, uj)
}

// CreateUser creates a new user resource
func (uc UserController) CreateUser(c *gin.Context) {

    var json user

    // This will infer what binder to use depending on the content-type header.
    c.Bind(&json) 

	// Stub an user to be populated from the body
	u := models.User{}

	user := createuser(json.Name, json.Address,json.Age)
    if user.Name == json.Name {
        content := gin.H{
            "result": "Success",
            "Name": user.Name,
            "Address": user.Address,
            "Age": user.Age
        }
    
        c.Header().Set("Content-Type", "application/json")
        c.JSON(201, content)
    } else {
        c.JSON(500, gin.H{"result": "An error occured"})
    }
	
}

// RemoveUser removes an existing user resource
func (uc UserController) RemoveUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove user
	if err := uc.session.DB(DB_NAME).C(DB_COLLECTION).RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}

func createuser(Name string, Address string,Age int) models.User {
    user := models.User{
        Name:      Name,
        Address:    Address,
        Age:	Age
    }
    // Write the user to mongo
	err := uc.session.DB(DB_NAME).C(DB_COLLECTION).Insert(&user)
    checkErr(err, "Insert failed")
    return user
}

func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}
