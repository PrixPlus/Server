# Prix

Introduction to the API

# Environment settings

To set the environment settings, please coppy the file `settings/dev.json` to folder `$GOPATH` an change its contents. The name of the file should be the same value in `GO_ENV` environment variable, or should be `dev.json` if it's not set.

# TESTs

To test the API, just enter `tests/api` folder and type `go test`


# API

#### [GET] /
Return the index.html dub page.


#### [GET] /api/ping
This route is used to check the service status.

- Return:
```javascript
{
  "message":"pong"
}
```


#### [POST] /api/login
This route is used to authenticate an user, it returns an valid token to make authenticated api calls.
This route uses `golang.org/x/crypto/bcrypt` that implements Provos and Mazières's bcrypt adaptive hashing algorithm,

- Request:
```javascript
{
  "email":"user@prix.plus",
  "password":"$2a$10$nLI0KOH0JGUqz3bjdtz6vOz6W/yo10suD.BC9Z8.rR9eBCUCzTOX."
}
```

- Return:
```javascript
{
  "expire":"2016-05-15T01:43:07-04:00",
  "token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjM3MTk0NjMsImlkIjoxfQ.5pd5X1WFk2t9PiuOb6T0D95mNJ5Fp7uxOBMOoiUh6adyewf64JSJZpv66y9EAjghPTvJ55bsxEhOxX-FcKr41Q"
}
```


#### [GET] /api/refresh_token
This route is used to refresh a valid token.

- Request Header: `Authorization: Bearer TOKEN`

- Return:
```javascript
{
  "expire":"2016-05-15T02:17:12-03:00",
  "token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjM3MTk0NjMsImlkIjoxfQ.5pd5X1WFk2t9PiuOb6T0D95mNJ5Fp7uxOBMOoiUh6adyewf64JSJZpv66y9EAjghPTvJ55bsxEhOxX-FcKr41Q"
}
```


#### [GET] /api/me
This route is used to retrieve the current user in the session.

- Request Header: `Authorization: Bearer TOKEN`

- Return:
```javascript
{
  "results":[
    {
      "id":1,
      "email":"user@prix.plus"
    }
  ]
}
```


#### [POST] /api/users
This route is used to create a new user.

- Request:
```javascript
{
  "email":"newuser@prix.plus",
  "password":"123456"
}
```

- Return Header: `Location: /api/users/2`
- Return:
```javascript
{
  "location":"/api/users/2",
  "results":[
    {
      "id":2,
      "email":"newuser@prix.plus"
    }
  ]
}
```


#### [PUT] /api/users/:id
This route is used to update users info. You can't use an `id` differente from yours.
Obs: **In this route all fields are optional**, so you can't resend your email

- Request Header: `Authorization: Bearer TOKEN`

- Request:
```javascript
{
  "id":1,
  "password":"123456",
  "email":"newuser@prix.plus"
}
```

- Return:
```javascript
{
  "results":[
    {
      "id":1,
      "password":"123456",
      "email":"newuser@prix.plus"
    }
  ]
}
```

# TESTS

Test user: Email `user@prix.plus` and Password `123456` 

#### Testing `/api/login` route

- Request:
```javascript
{
  "email":"user@prix.plus",
  "password":"$2a$10$nLI0KOH0JGUqz3bjdtz6vOz6W/yo10suD.BC9Z8.rR9eBCUCzTOX."
}
```