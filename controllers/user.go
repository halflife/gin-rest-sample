package controllers

import (
	"encoding/json"
	"log"
	"strconv"

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

func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}



// Get all Users
func (uc UserController) UsersList(c *gin.Context) {
	var results []models.User
    err := uc.session.DB(DB_NAME).C(DB_COLLECTION).Find(nil).All(&results)
	checkErr(err, "Users doesn't exist")
    content := gin.H{}
    for k, v := range results {
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
		c.AbortWithStatus(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Stub user
	u := models.User{}

	// Fetch user
	if err := uc.session.DB(DB_NAME).C(DB_COLLECTION).FindId(oid).One(&u); err != nil {
		checkErr(err, "Users doesn't exist")
		c.AbortWithStatus(403)

		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	c.Writer.Header().Set("Content-Type", "application/json")
    c.JSON(200, uj)
}

// CreateUser creates a new user resource
func (uc UserController) CreateUser(c *gin.Context) {

    var json models.User

    // This will infer what binder to use depending on the content-type header.
    c.Bind(&json) 

	// Stub an user to be populated from the body
	//u := models.User{}

	u := uc.create_user(json.Name, json.Gender,json.Age)
    if u.Name == json.Name {
        content := gin.H{
            "result": "Success",
            "Name": u.Name,
            "Gender": u.Gender,
            "Age": u.Age,
        }
    
        c.Writer.Header().Set("Content-Type", "application/json")
        c.JSON(201, content)
    } else {
        c.JSON(500, gin.H{"result": "An error occured"})
    }
	
}

// RemoveUser removes an existing user resource
func (uc UserController) RemoveUser(c *gin.Context) {
	// Grab id
	id := c.Params.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		c.AbortWithStatus(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove user
	if err := uc.session.DB(DB_NAME).C(DB_COLLECTION).RemoveId(oid); err != nil {
		checkErr(err,"Fail to Remove")
		c.AbortWithStatus(404)
		return
	}

	// Write status
	c.AbortWithStatus(200)
}


// RemoveUser removes an existing user resource
func (uc UserController) UpdateUser(c *gin.Context) {
	// Grab id
	id := c.Params.ByName("id")
    var json models.User

    // This will infer what binder to use depending on the content-type header.
    c.Bind(&json) 

	// Stub an user to be populated from the body
	//u := models.User{}

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		c.AbortWithStatus(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	u := uc.update_user(oid,json.Name, json.Gender,json.Age)
    if u.Name == json.Name {
        content := gin.H{
            "result": "Success",
            "Name": u.Name,
            "Gender": u.Gender,
            "Age": u.Age,
        }
    
        c.Writer.Header().Set("Content-Type", "application/json")
        c.JSON(201, content)
    } else {
        c.JSON(500, gin.H{"result": "An error occured"})
    }

	// Write status
	c.AbortWithStatus(200)
}

func (uc UserController) create_user(Name string, Gender string,Age int) models.User {
    user := models.User{
        Name:      Name,
        Gender:    Gender,
        Age:	Age,
    }
    // Write the user to mongo
	err := uc.session.DB(DB_NAME).C(DB_COLLECTION).Insert(&user)
    checkErr(err, "Insert failed")
    return user
}

func (uc UserController) update_user(Id bson.ObjectId,Name string, Gender string,Age int) models.User {

    user := models.User{
    	Id: 	Id,
        Name:      Name,
        Gender:    Gender,
        Age:	Age,
    }

    // Push a item to the Array in the Collection by Collection's ObjectId
    
    change := bson.M{"$push": bson.M{"sections": bson.M{"Name":user.Name,"Gender":user.Gender,"Age":user.Age}}}
    
    // Write the user to mongo
    if err := uc.session.DB(DB_NAME).C(DB_COLLECTION).Update(Id, bson.M{"$push": change}); err != nil {
		checkErr(err,"Update failed")
		//c.AbortWithStatus(404)
		//return
	}

    return user
}


