# gin-rest-sample
To get results 

run this in terminal

View all data

curl http://128.199.141.149:8000/users 
 
View buy ID

  curl http://128.199.141.149:8000/users/{id}
  
Add values

  curl -H "Content-Type: application/json" -d '{"name":"Amila","address":"23/2,Palavatha,Colombo 2"}' http://128.199.141.149:8000/usersadd
