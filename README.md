# gin-rest-sample

you need to import following packages

  go get github.com/gin-gonic/gin

  go get gopkg.in/mgo.v2

  go get gopkg.in/mgo.v2/bson

To run the server

  go run server.go


To get results 

run this in terminal

View all data

  curl http://128.199.141.149:8000/users 
 
View buy ID

  curl http://128.199.141.149:8000/users/{oid}
  
Add values

  curl -H "Content-Type: application/json" -d '{"name":"Amila","gender":"23/2,Palavatha,Colombo 2", "Age":34}' http://128.199.141.149:8000/usersadd

update values

  curl -H "Content-Type: application/json" -d '{"name":"Amila","gender":"23/2,Palavatha,pavithra Road,Colombo 2", "Age":34}' http://128.199.141.149:8000/usersupdate/{oid}

Remove a data
  curl http://128.199.141.149:8000/usersdelete/{oid}
