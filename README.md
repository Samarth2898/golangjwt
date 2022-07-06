# Golangjwt
An experiment to understand the implementation of JWT in Golang.

# How to run the project locally
- clone the repo
```
cd golangjwt
go run main.go
```

## Endpoints available
The server runs on port 9000 and can be accessed with the following end-points-
### user signup
```
http://localhost:9000/users/signup
```
- it is a POST request which takes the following as the payload.
```json
{
    "first_name":"abc",
    "last_name":"def",
    "password":"abc123%",
    "email":"abc@def.com",
    "phone":"1234567892",
    "user_type": "<user_type>" 
}
```
- user_type is string that takes the following values -
  - "ADMIN": user_type ADMIN has access to all routes(users, users/:id).
  - "USER": user_type USER has access to USER routes(users/:id)

### user login
```
http://localhost:9000/users/login
```
### get all users
```
http://localhost:9000/users
```
### get individual user
```
localhost:9000/users/:id
```
