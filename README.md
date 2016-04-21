# Prix

Introduction to the API

# Environment settings

To set the environment settings, please coppy the file `settings/dev.json` to folder `$GOPATH` an change its contents. The name of the file should be the same value in `GO_ENV` environment variable, or should be `dev.json` if it's not set.

# Testing

To test the API, just enter `tests/api` folder and type `go test`

For Application Test, please use: Test User: Email `test@test.com` and Password `123456`

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
  "email":"test@test.com",
  "password":"123456"
}
```

- Return:
```javascript
{
  "token": {
    "raw":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjM3MTk0NjMsImlkIjoxfQ.5pd5X1WFk2t9PiuOb6T0D95mNJ5Fp7uxOBMOoiUh6adyewf64JSJZpv66y9EAjghPTvJ55bsxEhOxX-FcKr41Q",
    "expire":"2016-05-15T01:43:07-04:00"
  }
}
```


#### [GET] /api/refresh_token
This route is used to refresh a valid token.

- Request Header: `Authorization: Bearer TOKEN`

- Return:
```javascript
{
  "token": {
    "raw":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjM3MTk0NjMsImlkIjoxfQ.5pd5X1WFk2t9PiuOb6T0D95mNJ5Fp7uxOBMOoiUh6adyewf64JSJZpv66y9EAjghPTvJ55bsxEhOxX-FcKr41Q",
    "expire":"2016-05-15T01:43:07-04:00"
  }
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
      "email":"test@test.com"
    }
  ]
}
```


#### [POST] /api/users
This route is used to create a new user.

- Request:
```javascript
{
  "email":"newtest@test.com",
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
      "email":"newtest@test.com"
    }
  ]
}
```


#### [PUT] /api/users/:id
This route is used to update users info. You can't use an `id` differente from yours.
Obs: **In this route all fields are optional**, so you can't resend your email
TODO: Not working for changing password yet

- Request Header: `Authorization: Bearer TOKEN`

- Request:
```javascript
{
  "id":1,
  "password":"123456", // TODO
  "email":"newtest@test.com"
}
```

- Return:
```javascript
{
  "results":[
    {
      "id":1,
      "email":"newtest@test.com"
    }
  ]
}
```