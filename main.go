package main

import (
    "github.com/gin-gonic/gin"
    "database/sql"
    "github.com/coopernurse/gorp"
    _ "github.com/mattn/go-sqlite3"
    "log"
    "time"
    "strconv"
)

var dbmap = initDb()

func main(){

    defer dbmap.Db.Close()

    router := gin.Default()
    router.GET("/users", usersList)
    router.POST("/usersadd", userPost)
    router.GET("/users/:id", usersDetail)
    router.Run(":8000")
}

type user struct {
    Id int64 `db:"user_id"`
    Created int64
    Name string
    Address string
}

func createuser(Name, body string) user {
    user := user{
        Created:    time.Now().UnixNano(),
        Name:      Name,
        Address:    body,
    }

    err := dbmap.Insert(&user)
    checkErr(err, "Insert failed")
    return user
}

func getuser(user_id int) user {
    user := user{}
    err := dbmap.SelectOne(&user, "select * from users where user_id=?", user_id)
    checkErr(err, "User doesn't exist")
    return user
}

func usersList(c *gin.Context) {
    var users []user
    _, err := dbmap.Select(&users, "select * from users order by user_id")
    checkErr(err, "Users doesn't exist")
    content := gin.H{}
    for k, v := range users {
        content[strconv.Itoa(k)] = v
    }
    c.JSON(200, content)
}

func usersDetail(c *gin.Context) {
    user_id := c.Params.ByName("id")
    a_id, _ := strconv.Atoi(user_id)
    user := getuser(a_id)
    content := gin.H{"Name": user.Name, "content": user.Address}
    c.JSON(200, content)
}

func userPost(c *gin.Context) {
    var json user

    c.Bind(&json) // This will infer what binder to use depending on the content-type header.
    user := createuser(json.Name, json.Address)
    if user.Name == json.Name {
        content := gin.H{
            "result": "Success",
            "Name": user.Name,
            "content": user.Address,
        }
        c.JSON(201, content)
    } else {
        c.JSON(500, gin.H{"result": "An error occured"})
    }
}

func initDb() *gorp.DbMap {
    db, err := sql.Open("sqlite3", "db.sqlite3")
    checkErr(err, "sql.Open failed")

    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

    dbmap.AddTableWithName(user{}, "users").SetKeys(true, "Id")

    err = dbmap.CreateTablesIfNotExists()
    checkErr(err, "Create tables failed")

    return dbmap
}

func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}